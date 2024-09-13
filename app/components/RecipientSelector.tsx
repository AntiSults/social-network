import React, { useState } from "react";

interface Recipient {
    id: number;
    name: string;
    type: "user" | "group";
}

interface RecipientSelectorProps {
    recipients: Recipient[];
    selectedRecipient: number | number[];
    setSelectedRecipient: React.Dispatch<React.SetStateAction<number | number[]>>;
}

const RecipientSelector: React.FC<RecipientSelectorProps> = ({
    recipients,
    selectedRecipient,
    setSelectedRecipient,
}) => {
    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedValue = e.target.value;
        const selectedID = parseInt(selectedValue, 10);

        if (!isNaN(selectedID)) {
            setSelectedRecipient(selectedID);
        }
    };

    return (
        <select
            value={selectedRecipient as number}
            onChange={handleChange}
            className="p-2 border border-gray-300 rounded-md"
        >
            <option value="" disabled>
                Select Recipient
            </option>
            {recipients.map((recipient) => (
                <option
                    key={recipient.id}
                    value={recipient.id}
                    style={{
                        color: recipient.type === "user" ? "blue" : "green",
                    }}
                >
                    {recipient.name} ({recipient.type})
                </option>
            ))}
        </select>
    );
};

export default RecipientSelector;
