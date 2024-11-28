import React from "react";
import { Form, Input, Button } from "antd";

// Login Form
export const LoginForm: React.FC<{ onSubmit: (values: { email: string; password: string }) => void }> = ({ onSubmit }) => {
  const [form] = Form.useForm();
  const handleSubmit = () => form.validateFields().then(onSubmit);

  return (
    <Form form={form} layout="vertical">
      <Form.Item label="Email" name="email" rules={[{ required: true, type: "email", message: "Введите корректный email" }]}>
        <Input placeholder="Email" />
      </Form.Item>
      <Form.Item label="Password" name="password" rules={[{ required: true, message: "Введите пароль" }]}>
        <Input.Password placeholder="Password" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" onClick={handleSubmit}>Войти</Button>
      </Form.Item>
    </Form>
  );
};

// Register Form
export const RegisterForm: React.FC<{ onSubmit: (values: { email: string; password: string; username: string }) => void }> = ({ onSubmit }) => {
  const [form] = Form.useForm();
  const handleSubmit = () => form.validateFields().then(onSubmit);

  return (
    <Form form={form} layout="vertical">
      <Form.Item label="Username" name="username" rules={[{ required: true, message: "Введите имя пользователя" }]}>
        <Input placeholder="Username" />
      </Form.Item>
      <Form.Item label="Email" name="email" rules={[{ required: true, type: "email", message: "Введите корректный email" }]}>
        <Input placeholder="Email" />
      </Form.Item>
      <Form.Item label="Password" name="password" rules={[{ required: true, message: "Введите пароль" }]}>
        <Input.Password placeholder="Password" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" onClick={handleSubmit}>Зарегистрироваться</Button>
      </Form.Item>
    </Form>
  );
};
