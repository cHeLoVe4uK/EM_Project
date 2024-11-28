import React from "react";
import { LoginForm } from "../components/Form";
import { message } from "antd";

const LoginPage: React.FC = () => {
  const handleLogin = async (values: { email: string; password: string }) => {
    try {
      const response = await fetch("/api/v1/user/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(values),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.msg || "Login failed");
      }

      message.success("Login successful");
      // Navigate to ChatPage after successful login
    } catch (err: any) {
      message.error(err.message);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: "auto", padding: "1rem"}}>
      <h2>Login</h2>
      <LoginForm onSubmit={handleLogin} />
    </div>
  );
};

export default LoginPage;
