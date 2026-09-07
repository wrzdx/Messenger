// A deployment smoke check. Creates two temporary accounts and an empty direct
// chat, removes the test message, and anonymizes the accounts through the API.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type account struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	Token string `json:"access_token"`
}

func request(ctx context.Context, base, method, path, token string, body any, expected int, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequestWithContext(ctx, method, base+"/api/v1"+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != expected {
		return fmt.Errorf("%s %s: expected %d, got %d", method, path, expected, resp.StatusCode)
	}
	if out == nil {
		_, err = io.Copy(io.Discard, resp.Body)
		return err
	}
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	return json.Unmarshal(envelope.Data, out)
}

func run(base string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	var accounts []account
	defer func() {
		cleanup, stop := context.WithTimeout(context.Background(), 15*time.Second)
		defer stop()
		for _, a := range accounts {
			if err := request(cleanup, base, "DELETE", "/users/me", a.Token, nil, 204, nil); err != nil {
				fmt.Fprintln(os.Stderr, "account cleanup:", err)
			}
		}
	}()
	for range 2 {
		var a account
		err := request(ctx, base, "POST", "/auth/register", "", map[string]string{
			"username":   "smoke_" + strings.ReplaceAll(uuid.NewString(), "-", "")[:16],
			"first_name": "Deployment check", "password": uuid.NewString(),
		}, 201, &a)
		if err != nil {
			return err
		}
		accounts = append(accounts, a)
	}
	fmt.Println("registration: OK (two accounts)")
	wsURL := "ws" + strings.TrimPrefix(base, "http") + "/api/v1/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": []string{base}}})
	if err != nil {
		return fmt.Errorf("websocket handshake: %w", err)
	}
	defer conn.CloseNow()
	if err = wsjson.Write(ctx, conn, map[string]string{"type": "authenticate", "access_token": accounts[1].Token}); err != nil {
		return err
	}
	var event struct {
		Type string `json:"type"`
		Data struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		} `json:"data"`
	}
	if err = wsjson.Read(ctx, conn, &event); err != nil {
		return err
	}
	if event.Type != "authenticated" {
		return fmt.Errorf("unexpected auth response: %s", event.Type)
	}
	fmt.Println("WebSocket upgrade and authentication: OK")
	var chat struct {
		ID string `json:"id"`
	}
	if err = request(ctx, base, "POST", "/chats/directs", accounts[0].Token, map[string]string{"peer_id": accounts[1].User.ID}, 201, &chat); err != nil {
		return err
	}
	path := "/chats/" + chat.ID + "/messages"
	var msg struct {
		ID string `json:"id"`
	}
	if err = request(ctx, base, "POST", path+"/", accounts[0].Token, map[string]string{"client_message_id": uuid.NewString(), "content": "Deployment smoke check"}, 201, &msg); err != nil {
		return err
	}
	deleted := false
	defer func() {
		if !deleted {
			cleanup, stop := context.WithTimeout(context.Background(), 10*time.Second)
			defer stop()
			if err := request(cleanup, base, "DELETE", path+"/"+msg.ID, accounts[0].Token, nil, 204, nil); err != nil {
				fmt.Fprintln(os.Stderr, "message cleanup:", err)
			}
		}
	}()
	for _, kind := range []string{"message_created", "message_edited", "message_deleted"} {
		switch kind {
		case "message_edited":
			err = request(ctx, base, "PATCH", path+"/"+msg.ID, accounts[0].Token, map[string]string{"content": "Updated deployment check"}, 200, nil)
		case "message_deleted":
			err = request(ctx, base, "PUT", path+"/read", accounts[1].Token, map[string]string{"message_id": msg.ID}, 204, nil)
			if err == nil {
				err = request(ctx, base, "DELETE", path+"/"+msg.ID, accounts[0].Token, nil, 204, nil)
				deleted = err == nil
			}
		}
		if err != nil {
			return err
		}
		if err = wsjson.Read(ctx, conn, &event); err != nil {
			return err
		}
		if event.Type != kind || event.Data.ID != msg.ID {
			return fmt.Errorf("unexpected event while waiting for %s", kind)
		}
		fmt.Println(kind + ": OK")
	}
	var page struct {
		Messages []json.RawMessage `json:"messages"`
	}
	if err = request(ctx, base, "GET", path+"/", accounts[1].Token, nil, 200, &page); err != nil {
		return err
	}
	if len(page.Messages) != 0 {
		return fmt.Errorf("deleted message remains in history")
	}
	fmt.Println("read marker and empty history after deletion: OK")
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: smoke https://messenger.example.com")
		os.Exit(2)
	}
	if err := run(strings.TrimRight(os.Args[1], "/")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
