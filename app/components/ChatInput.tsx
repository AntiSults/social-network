import React, { useState } from "react";
import FieldInput from "./FieldInput";
import Button from "./Button";
import EmojiPicker, { EmojiClickData } from "emoji-picker-react";
import { getOptionStyle } from "@/app/utils/getOptionStyle";
import { Recipient } from "@/app/utils/types";

interface ChatInputProps {
    message: string;
    setMessage: React.Dispatch<React.SetStateAction<string>>;
    onSubmit: (e: React.FormEvent) => void;
    recipients: Recipient[];
    selectedRecipient: number | null;
    setSelectedRecipient: React.Dispatch<React.SetStateAction<number | null>>;
    groupId: string | null;
}

const ChatInput: React.FC<ChatInputProps> = ({
    message,
    setMessage,
    onSubmit,
    recipients,
    selectedRecipient,
    setSelectedRecipient,
    groupId,
}) => {
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);

    const onEmojiClick = (emojiObject: EmojiClickData) => {
        setMessage(message + emojiObject.emoji);
        setShowEmojiPicker(false);
    };

    return (
        <form onSubmit={onSubmit} className="flex flex-col space-y-2">
            {/* Recipient Selector */}
            <select
                value={selectedRecipient ?? ""}  // Single recipient, default empty string
                onChange={(e) => {
                    const selectedID = parseInt(e.target.value);
                    setSelectedRecipient(selectedID);  // Directly set the selected recipient
                }}
                className="p-2 border border-gray-300 rounded-md"
            >
                <option value="" disabled>Select a recipient</option> {/* Default option */}
                {recipients.map((recipient) => (
                    <option key={recipient.id} value={recipient.id} style={getOptionStyle(recipient.type)}>
                        {recipient.name}
                    </option>
                ))}
            </select>

            <div className="relative flex items-center">
                <FieldInput
                    name="Text:"
                    type="text"
                    placeholder="Push your imagination"
                    required={true}
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    className="p-2 border border-gray-300 rounded-md flex-grow"
                />
                <button
                    type="button"
                    onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                    className="ml-2 p-2 bg-gray-200 rounded-md"
                >
                    😊
                </button>
            </div>

            {/* Emoji Picker */}
            {showEmojiPicker && (
                <div className="absolute z-10">
                    <EmojiPicker onEmojiClick={onEmojiClick} />
                </div>
            )}

            <Button type="submit" name="Submit Message" className="p-2 bg-blue-500 text-white rounded-md" />
        </form>
    );
};

export default ChatInput;
