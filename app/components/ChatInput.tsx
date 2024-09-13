import React, { useState } from "react";
import FieldInput from "./FieldInput";
import Button from "./Button";
import EmojiPicker, { EmojiClickData } from "emoji-picker-react"; // Import EmojiClickData

interface Recipient {
    id: number;
    name: string;
    type: "user" | "group";
}

interface ChatInputProps {
    message: string;
    setMessage: React.Dispatch<React.SetStateAction<string>>;
    onSubmit: (e: React.FormEvent) => void;
    recipients: Recipient[];
    selectedRecipient: number | number[];
    setSelectedRecipient: React.Dispatch<React.SetStateAction<number | number[]>>;
}

const ChatInput: React.FC<ChatInputProps> = ({
    message,
    setMessage,
    onSubmit,
    recipients,
    selectedRecipient,
    setSelectedRecipient,
}) => {
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);

    // Correct the function to use the right signature
    const onEmojiClick = (emojiObject: EmojiClickData) => {
        setMessage(message + emojiObject.emoji); // Append the selected emoji to the message
        setShowEmojiPicker(false); // Close the emoji picker after selecting an emoji
    };

    const getOptionStyle = (type: "user" | "group") => {
        return type === "user" ? { color: "blue" } : { color: "green" };
    };

    return (
        <form onSubmit={onSubmit} className="flex flex-col space-y-2">
            {/* Recipient Selector */}
            <select
                value={Array.isArray(selectedRecipient) ? selectedRecipient[0] : selectedRecipient}
                onChange={(e) =>
                    setSelectedRecipient(
                        recipients.find((r) => r.id === parseInt(e.target.value))?.id || 0
                    )
                }
                className="p-2 border border-gray-300 rounded-md"
            >
                {recipients.map((recipient) => (
                    <option key={recipient.id} value={recipient.id} style={getOptionStyle(recipient.type)}>
                        {recipient.name} ({recipient.type})
                    </option>
                ))}
            </select>

            {/* Text Input */}
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
                {/* Emoji Button */}
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

            {/* Submit Button */}
            <Button type="submit" name="Submit Message" className="p-2 bg-blue-500 text-white rounded-md" />
        </form>
    );
};

export default ChatInput;