import React, { useState, useEffect, useRef } from "react";
import { Button, Input, message, Modal, Layout, List } from "antd";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { createChat, connectToChat } from "../api/api";

const { Sider, Content } = Layout;

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
    if (selectedChat) {
      const socket = connectToChat({ chat_id: selectedChat });
      socketRef.current = socket;
  
      socket.onopen = () => {
        console.log("WebSocket connected");
        message.success("Connected to chat!");
      };
  
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        setMessages((prev) => [...prev, data]);
      };
  
      socket.onerror = () => {
        message.error("WebSocket connection error. Check the server.");
      };
  
      socket.onclose = () => {
        console.log("WebSocket connection closed");
      };
  
      return () => {
        if (socketRef.current) {
          socketRef.current.close();
          socketRef.current = null;
        }
      };
    }
  }, [selectedChat]);  

  const handleSelectChat = async (chatId: string) => {
    setSelectedChat(chatId);

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

  const handleSendMessage = () => {
    if (!selectedChat || !socketRef.current) {
      message.warning("Please select a chat or check your connection.");
      return;
    }

    const messageData = { text: newMessage, sender: "You" };
    socketRef.current.send(JSON.stringify(messageData));

    setMessages((prev) => [
      ...prev,
      { id: Date.now().toString(), ...messageData },
    ]);

    setNewMessage("");
  };

  const handleCreateChat = async () => {
    if (!newChatName.trim()) {
      message.warning("Chat name cannot be empty.");
      return;
    }

    try {
      const newChat = await createChat({ name: newChatName });
      setChats((prev) => [...prev, { id: newChat.id, name: newChatName }]);
      setNewChatName("");
      setIsModalVisible(false);
      message.success("Chat created successfully.");
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
    <Layout style={{ height: "100vh" }}>
      <Sider
        theme="light"
        width={250}
        style={{ borderRight: "1px solid #ccc", overflowY: "auto" }}
      >
        <List
          header={<h3>Chats</h3>}
          dataSource={chats}
          renderItem={(chat) => (
            <List.Item
              style={{
                cursor: "pointer",
                padding: "10px",
                background: chat.id === selectedChat ? "#f0f0f0" : "inherit",
              }}
              onClick={() => handleSelectChat(chat.id)}
            >
              {chat.name}
            </List.Item>
          )}
        />
        <Button type="primary" onClick={() => setIsModalVisible(true)} block>
          New Chat
        </Button>
        <Button type="primary" danger onClick={handleLogout} block>
          Logout
        </Button>
      </Sider>

      <Content style={{ display: "flex", flexDirection: "column" }}>
        <div style={{ flex: 1, overflowY: "auto", padding: "1rem" }}>
          {selectedChat ? (
            <List
              header={<h3>Messages</h3>}
              dataSource={messages}
              renderItem={(message) => (
                <List.Item>
                  <strong>{message.sender}:</strong> {message.text}
                </List.Item>
              )}
            />
          ) : (
            <p>Select a chat to start messaging.</p>
          )}
        </div>

        <div style={{ padding: "1rem", display: "flex", gap: "10px" }}>
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
      </Content>

      <Modal
        title="Create New Chat"
        open={isModalVisible}
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
    </Layout>
  );
};

export default ChatPage;
