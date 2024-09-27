
import { User as SearchResult } from "@/app/utils/types";

export const searchUsers = async (query: string): Promise<SearchResult[]> => {
    try {
        const response = await fetch(`http://localhost:8080/search?query=${encodeURIComponent(query)}`);
        if (response.ok) {
            return await response.json();
        } else {
            console.error("Search failed");
            return [];
        }
    } catch (error) {
        console.error("Error during search:", error);
        return [];
    }
};