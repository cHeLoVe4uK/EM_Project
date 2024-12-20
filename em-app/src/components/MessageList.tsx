import React from "react";
import { List, Typography } from "antd";

interface Message {
  id: string;
  text: string;
  sender: string;
}

interface MessageListProps {
  messages: Message[];
}

const MessageList: React.FC<MessageListProps> = ({ messages }) => {
  return (
    <List
      dataSource={messages}
      renderItem={(message) => (
        <List.Item>
          <Typography.Text strong>{message.sender}: </Typography.Text>
          <Typography.Text>{message.text}</Typography.Text>
        </List.Item>
      )}
    />
  );
};

export default MessageList;
