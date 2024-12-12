import React, { useState, useEffect, useRef } from "react";
import { Button, Input, message, Modal, Layout, List, Typography } from "antd";
import { PlusOutlined, LogoutOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { createChat, connectToChat } from "../api/api";
import { Message as MessageType } from "../api/types";

const { Sider, Content } = Layout;
const { Title } = Typography;

interface Chat {
  id: string;
  name: string;
}

const ChatPage: React.FC = () => {
  const [chats, setChats] = useState<Chat[]>([]);
  const [messages, setMessages] = useState<MessageType[]>([]);
  const [selectedChat, setSelectedChat] = useState<string | null>(null);
  const [newMessage, setNewMessage] = useState("");
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [newChatName, setNewChatName] = useState("");
  const [isLogoutModalVisible, setIsLogoutModalVisible] = useState(false);
  const navigate = useNavigate();
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const fetchChats = async () => {
      try {
        const token = Cookies.get("access_token");
        if (!token) {
          throw new Error("No access token found");
        }
    
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/chats`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
    
        if (!response.ok) {
          if (response.status === 401) {
            throw new Error("Unauthorized: Please login again.");
          }
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
      const socket = connectToChat(selectedChat);
      socketRef.current = socket;
  
      socket.onopen = () => {
        console.log("WebSocket connected");
        message.success("Connected to chat!");
      };
  
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
  
        if (data && data.content && data.created_at) {
          const formattedDate = new Date(data.created_at).toLocaleString();
  
          setMessages((prev) => [
            ...prev,
            { ...data, created_at: formattedDate, is_edited: false },
          ]);
        }
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
      const token = Cookies.get("access_token");
      if (!token) {
        throw new Error("No access token found");
      }
  
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/chats/${chatId}/messages`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
  
      if (!response.ok) {
        if (response.status === 401) {
          throw new Error("Unauthorized: Please login again.");
        }
        throw new Error("Failed to fetch messages");
      }
  
      const data = await response.json();
      setMessages(data);
    } catch (err: any) {
      message.error(err.message);
    }
  };
  

  const handleSendMessage = async () => {
    if (!newMessage.trim()) {
      message.warning("Message cannot be empty.");
      return;
    }
  
    if (!selectedChat) {
      message.warning("Please select a chat first.");
      return;
    }
  
    try {
      const token = Cookies.get("access_token");
      if (!token) {
        throw new Error("No access token found");
      }
  
      const messageData = {
        content: newMessage,
        created_at: new Date().toISOString(),
      };
  
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/chats/${selectedChat}/messages`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(messageData),
      });
  
      if (!response.ok) {
        if (response.status === 401) {
          throw new Error("Unauthorized: Please login again.");
        }
        throw new Error("Failed to send message");
      }
  
      const sentMessage = await response.json();
      setMessages((prev) => [...prev, sentMessage]);
      setNewMessage("");
    } catch (err: any) {
      message.error(err.message);
    }
  };
  

  const handleCreateChat = async () => {
    if (!newChatName.trim()) {
      message.warning("Chat name cannot be empty.");
      return;
    }
  
    try {
      const token = Cookies.get("access_token");
      if (!token) {
        throw new Error("No access token found");
      }
  
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/chats`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name: newChatName }),
      });
  
      if (!response.ok) {
        if (response.status === 401) {
          throw new Error("Unauthorized: Please login again.");
        }
        throw new Error("Failed to create chat");
      }
  
      const newChat = await response.json();
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

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  return (
    <Layout style={{ height: "100vh", background: "#f0f2f5" }}>
      <Sider
        theme="light"
        width={300}
        style={{
          borderRight: "1px solid #ccc",
          overflowY: "auto",
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
        }}
      >
        <div>
          <Title level={3} style={{ textAlign: "center", marginBottom: "20px", padding: "20px", borderBottom: "1px solid #ccc" }}>
            Chats
          </Title>
          <List
            dataSource={chats}
            renderItem={(chat) => (
              <List.Item
                style={{
                  cursor: "pointer",
                  padding: "10px",
                  background: chat.id === selectedChat ? "#e6f7ff" : "#fff",
                  borderRadius: "5px",
                  marginBottom: "10px",
                }}
                onClick={() => handleSelectChat(chat.id)}
              >
                {chat.name}
              </List.Item>
            )}
          />
          
          <Button
            icon={<PlusOutlined />}
            onClick={() => setIsModalVisible(true)}
            style={{
              width: "80%",
              background: "#ffffff",
              borderRadius: "12",
              left: 30,
              padding: "20px",
            }}
          />
        </div>
        <Button
          icon={<LogoutOutlined />}
          type="primary"
          danger
          onClick={() => setIsLogoutModalVisible(true)}
          style={{ position: 'absolute', bottom: 20, left: 20 }}
        >
          Logout
        </Button>
      </Sider>

      <Content style={{ padding: "20px", display: "flex", flexDirection: "column" }}>
        <div
          style={{
            flex: 1,
            overflowY: "auto",
            border: "1px solid #ccc",
            borderRadius: "10px",
            padding: "20px",
            background: "#fff",
          }}
        >
          {selectedChat ? (
            <List
              dataSource={messages}
              renderItem={(message) => (
                <List.Item
                  style={{
                    padding: "10px 15px",
                    marginBottom: "10px",
                    borderRadius: "12px",
                    background: message.author === "You" ? "#d9f7be" : "#f0f0f0",
                    alignSelf: message.author === "You" ? "flex-end" : "flex-start",
                    maxWidth: "70%",
                    boxShadow: "0 2px 5px rgba(0, 0, 0, 0.1)",
                  }}
                >
                  <div>
                    <Typography.Text
                      style={{
                        fontWeight: message.author === "You" ? "bold" : "normal",
                        marginBottom: "5px",
                        display: "block",
                      }}
                    >
                      {message.author}
                    </Typography.Text>
                    <Typography.Text>{message.content}</Typography.Text>
                    <Typography.Text
                      style={{
                        fontSize: "12px",
                        color: "#888",
                        marginTop: "5px",
                        display: "block",
                      }}
                    >
                      {formatDate(message.created_at)}
                    </Typography.Text>
                  </div>
                </List.Item>
              )}
            />
          ) : (
            <p style={{ textAlign: "center", color: "#888" }}>
              Select a chat to start messaging.
            </p>
          )}
        </div>

        <div style={{ marginTop: "10px", display: "flex", gap: "10px" }}>
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

      <Modal
        title="Confirm Logout"
        open={isLogoutModalVisible}
        onOk={handleLogout}
        onCancel={() => setIsLogoutModalVisible(false)}
        okText="Logout"
        cancelText="Cancel"
      >
        Are you sure you want to logout?
      </Modal>
    </Layout>
  );
};

export default ChatPage;
