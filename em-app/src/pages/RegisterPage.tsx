import React from "react";
import { RegisterForm } from "../components/Form";
import { message } from "antd";

const RegisterPage: React.FC = () => {
  const handleRegister = async (values: { email: string; password: string; username: string }) => {
    try {
      const response = await fetch("/api/v1/user/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(values),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.msg || "Registration failed");
      }

      message.success("Registration successful");
      // Navigate to LoginPage after successful registration
    } catch (err: any) {
      message.error(err.message);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: "auto", padding: "1rem" }}>
      <h2>Register</h2>
      <RegisterForm onSubmit={handleRegister} />
    </div>
  );
};

export default RegisterPage;
