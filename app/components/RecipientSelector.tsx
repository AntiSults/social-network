import React from "react";
import { getOptionStyle } from "../utils/getOptionStyle";

interface Recipient {
    id: number;
    name: string;
    type: "user" | "group";
}

interface RecipientSelectorProps {
    recipients: Recipient[];
    selectedRecipient: number[];
    setSelectedRecipient: React.Dispatch<React.SetStateAction<number[]>>;
}

const RecipientSelector: React.FC<RecipientSelectorProps> = ({
    recipients,
    selectedRecipient,
    setSelectedRecipient,
}) => {
    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedValue = parseInt(e.target.value, 10);
        setSelectedRecipient([selectedValue]); // Wrap in array
    };

    return (
        <select
            value={selectedRecipient.length > 0 ? selectedRecipient[0] : ""}
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
                    {recipient.name} ({recipient.type})
                </option>
            ))}
        </select>
    );
};

export default RecipientSelector;
