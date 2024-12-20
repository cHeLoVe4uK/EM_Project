import React from "react";
import { message, Button, Form, Input } from "antd";
import { useNavigate } from "react-router-dom";

const LoginPage: React.FC = () => {
  const navigate = useNavigate();

  const handleLogin = async (values: { email: string; password: string }) => {
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/users/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(values),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.msg || "Login failed");
      }

      const data = await response.json();
      document.cookie = `access_token=${data.access_token}`;
      document.cookie = `refresh_token=${data.refresh_token}`;
      message.success("Login successful");
      navigate("/chat");
    } catch (err: any) {
      message.error(err.message);
    }
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        height: "100vh",
        backgroundColor: "#041528",
      }}
    >
      <h1
        style={{
          textAlign: "center",
          fontSize: "5em",
          margin: "0px",
          marginTop: "70px",
        }}
      >
        <a href="/" style={{ color: "#fdfdfd", textDecoration: "none" }}>
          Chat.
        </a>
      </h1>
      <h3
        style={{
          color: "#aaaaaa",
          textAlign: "center",
          fontSize: "1.6em",
          margin: "0px",
        }}
      >
        -Effective Mobile-
      </h3>
      <Form
        onFinish={handleLogin}
        style={{
          maxWidth: "300px",
          width: "100%",
          paddingRight: "25px",
          paddingTop: "35px",
          paddingBottom: "20px",
          paddingLeft: "25px",
          border: "1px solid #fdfdfd",
          borderRadius: "15px",
          backgroundColor: "#dfdfdf",
          marginTop: "80px",
        }}
      >
        <Form.Item
          name="email"
          rules={[{ required: true, message: "Please input your email!" }]}
        >
          <Input placeholder="Email" />
        </Form.Item>
        <Form.Item
          name="password"
          rules={[{ required: true, message: "Please input your password!" }]}
        >
          <Input.Password placeholder="Password" />
        </Form.Item>
        <Form.Item>
          <div style={{ display: "flex", justifyContent: "space-between" }}>
            <Button type="primary" htmlType="submit">
              Login
            </Button>
            <Button type="default" onClick={() => navigate("/register")}>
              Register
            </Button>
          </div>
        </Form.Item>
      </Form>
    </div>
  );
};

export default LoginPage;
