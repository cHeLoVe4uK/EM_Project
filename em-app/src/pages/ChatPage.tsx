import React, { useState, useEffect, useRef } from "react";
import { Button, Input, message, Modal } from "antd";
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
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [newChatName, setNewChatName] = useState("");
  const navigate = useNavigate();
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const fetchChats = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/chats`);
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

  useEffect(() => {
    const socket = new WebSocket(`${import.meta.env.VITE_WS_BASE_URL}/ws`);
    socketRef.current = socket;

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.chat_id === selectedChat) {
        setMessages((prev) => [...prev, data]);
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
      message.error("WebSocket connection error");
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };

    return () => {
      socket.close();
    };
  }, [selectedChat]);

  const handleSelectChat = (chatId: string) => {
    setSelectedChat(chatId);

    const fetchMessages = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/chats/${chatId}/messages`);
        if (!response.ok) {
          throw new Error("Failed to fetch messages");
        }
        const data = await response.json();
        setMessages(data);
      } catch (err: any) {
        message.error(err.message);
      }
    };

    fetchMessages();
  };

  const handleSendMessage = () => {
    if (!selectedChat || !socketRef.current) {
      message.warning("Please select a chat or check your connection");
      return;
    }

    const messageData = {
      chat_id: selectedChat,
      text: newMessage,
      sender: "You",
    };

    socketRef.current.send(JSON.stringify(messageData));

    setMessages((prev) => [
      ...prev,
      { id: Date.now().toString(), ...messageData },
    ]);

    setNewMessage("");
  };

  const handleCreateChat = async () => {
    if (!newChatName.trim()) {
      message.warning("Chat name cannot be empty");
      return;
    }

    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/chats`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: newChatName }),
      });

      if (!response.ok) {
        throw new Error("Failed to create chat");
      }

      const newChat = await response.json();
      setChats((prev) => [...prev, newChat]);
      message.success("Chat created successfully");
      setNewChatName("");
      setIsModalVisible(false);
    } catch (err: any) {
      message.error(err.message);
    }
  };

  const handleLogout = () => {
    Cookies.remove("access_token");
    Cookies.remove("refresh_token");
    if (socketRef.current) {
      socketRef.current.close();
    }
    navigate("/login");
  };

  return (
    <div style={{ display: "flex", height: "100vh" }}>
      <div
        style={{
          flex: 1,
          borderRight: "1px solid #ccc",
          overflowY: "auto",
          display: "flex",
          flexDirection: "column",
          position: "relative",
        }}
      >
        <ChatList chats={chats} onSelectChat={handleSelectChat} />
        <Button
          type="primary"
          style={{
            position: "absolute",
            bottom: "70px",
            left: "10px",
            width: "90%",
          }}
          onClick={() => setIsModalVisible(true)}
        >
          New Chat
        </Button>
        <Button
          type="primary"
          danger
          style={{
            position: "absolute",
            bottom: "20px",
            left: "10px",
            width: "90%",
          }}
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
      <Modal
        title="Create New Chat"
        visible={isModalVisible}
        onOk={handleCreateChat}
        onCancel={() => setIsModalVisible(false)}
        okText="Create"
        cancelText="Cancel"
      >
        <Input
          value={newChatName}
          onChange={(e) => setNewChatName(e.target.value)}
          placeholder="Enter chat name"
        />
      </Modal>
    </div>
  );
};

export default ChatPage;
