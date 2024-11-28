import React, { useState, useEffect } from "react";
import { Button, Input, message } from "antd";
import ChatList from "../components/ChatList";
import MessageList from "../components/MessageList";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";

interface Chat {
  id: string;
  name: string;
}

interface Message {
  id: string;
  text: string;
  sender: string;
}

const ChatPage: React.FC = () => {
  const [chats, setChats] = useState<Chat[]>([]);
  const [messages, setMessages] = useState<Message[]>([]);
  const [selectedChat, setSelectedChat] = useState<string | null>(null);
  const [newMessage, setNewMessage] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const fetchChats = async () => {
      try {
        const response = await fetch("/api/v1/chat");
        if (!response.ok) {
          throw new Error("Failed to fetch chats");
        }
        const data = await response.json();
        setChats(data);
      } catch (err: any) {
        message.error(err.message);
      }
    };

    fetchChats();
  }, []);

  const handleSelectChat = (chatId: string) => {
    setSelectedChat(chatId);
    // Fetch messages for the selected chat
  };

  const handleSendMessage = async () => {
    if (!selectedChat) {
      message.warning("Please select a chat");
      return;
    }

    try {
      const response = await fetch("/api/v1/message", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ chat_id: selectedChat, text: newMessage }),
      });

      if (!response.ok) {
        throw new Error("Failed to send message");
      }

      setMessages((prev) => [...prev, { id: Date.now().toString(), text: newMessage, sender: "You" }]);
      setNewMessage("");
    } catch (err: any) {
      message.error(err.message);
    }
  };

  const handleLogout = () => {
    // Удаляем токены из cookies
    Cookies.remove("access_token");
    Cookies.remove("refresh_token");
    // Редирект на страницу авторизации
    navigate("/login");
  };

  return (
    <div style={{ display: "flex", height: "100vh" }}>
      <div style={{ flex: 1, borderRight: "1px solid #ccc", overflowY: "auto", display: "flex", flexDirection: "column" }}>
        <ChatList chats={chats} onSelectChat={handleSelectChat} />
        <Button
          type="primary"
          danger
          style={{ position: "absolute", bottom: "20px", left: "10px", width: "50px" }}
          onClick={handleLogout}
        >
          Logout
        </Button>
      </div>
      <div style={{ flex: 3, display: "flex", flexDirection: "column" }}>
        <div style={{ flex: 1, overflowY: "auto", padding: "1rem" }}>
          <MessageList messages={messages} />
        </div>
        <div style={{ display: "flex", padding: "1rem" }}>
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type your message..."
            onPressEnter={handleSendMessage}
          />
          <Button type="primary" onClick={handleSendMessage}>
            Send
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ChatPage;
