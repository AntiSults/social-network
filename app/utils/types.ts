export interface User {
    ID: number;
    email: string;
    firstName: string;
    lastName: string;
    dob: string;
    nickname?: string;
    aboutMe?: string;
    avatarPath?: string | null;
    profileVisibility?: "public" | "private";
}

export interface Group {
    id: number;
    name: string;
    description: string;
    creator_id: number;
    created_at: string;
    members: number[];
}