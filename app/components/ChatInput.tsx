import React from "react";
import FieldInput from "./FieldInput";
import Button from "./Button";

interface ChatInputProps {
    message: string;
    setMessage: React.Dispatch<React.SetStateAction<string>>;
    onSubmit: (e: React.FormEvent) => void;
}

const ChatInput: React.FC<ChatInputProps> = ({ message, setMessage, onSubmit }) => {
    return (
        <form onSubmit={onSubmit} className="flex items-center space-x-2">
            <FieldInput
                name="Text:"
                type="text"
                placeholder="Push your imagination"
                required={true}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                className="flex-grow p-2 border border-gray-300 rounded-md"
            />
            <Button type="submit" name="Submit Message" className="p-2 bg-blue-500 text-white rounded-md" />
        </form>
    );
};

export default ChatInput;
