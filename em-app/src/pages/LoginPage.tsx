import React from "react";
import { LoginForm } from "../components/Form";
import { message } from "antd";
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
    <div style={{ maxWidth: 400, margin: "auto", padding: "1rem" }}>
      <h2>Login</h2>
      <LoginForm onSubmit={handleLogin} />
    </div>
  );
};

export default LoginPage;
