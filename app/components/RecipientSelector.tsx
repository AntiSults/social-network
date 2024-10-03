import React from "react";
import { getOptionStyle } from "@/app/utils/getOptionStyle";
import { Recipient } from "@/app/utils/types";

interface RecipientSelectorProps {
    recipients: Recipient[];
    selectedRecipient: number | null;
    setSelectedRecipient: React.Dispatch<React.SetStateAction<number | null>>;
}

const RecipientSelector: React.FC<RecipientSelectorProps> = ({
    recipients,
    selectedRecipient,
    setSelectedRecipient,
}) => {
    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedValue = parseInt(e.target.value, 10);
        setSelectedRecipient(selectedValue); // Directly set the selected recipient
    };

    return (
        <select
            value={selectedRecipient !== null ? selectedRecipient : ""}
            onChange={handleChange}
            className="p-2 border border-gray-300 rounded-md"
        >
            <option value="" disabled>Select Recipient</option>
            {recipients.map((recipient) => (
                <option
                    key={recipient.id}
                    value={recipient.id}
                    style={getOptionStyle(recipient.type)}
                >
                    {recipient.name}
                </option>
            ))}
        </select>
    );
};

export default RecipientSelector;
