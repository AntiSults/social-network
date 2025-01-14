export interface User {
    ID: number;
    email: string;
    firstName: string;
    lastName: string;
    dob: string;
    nickname?: string;
    aboutMe?: string;
    avatarPath?: string | null;
    profileVisibility: "public" | "private";
}

export interface Group {
    id: number;
    name: string;
    description: string;
    creator_id: number;
    created_at: string;
    members: number[];
}

export interface GroupJoinRequest {
    id: number;
    group_id: number;
    user_id: number;
    status: string;
    group_name: string;
    first_name: string;
    last_name: string;
}

export interface Recipient {
    id: number;
    name: string;
    type: "user" | "group";
}

export interface Message {
    id?: number;
    content: string;
    fromUserID: number;
    toUserID: number | number[];
    created: string;
    groupID: number | null;
}

export interface Post {
    id: number;
    content: string;
    created_at: string;
    author_first_name: string;
    author_last_name: string;
    files: string | null;
    group_name: string | null;
    privacy: string | null;
    allowedUserIds?: number[];
}

export interface Comment {
    id: number;
    postID: number;
    userID: number;
    content: string;
    file?: string;
    created_at: string;
    author_first_name: string;
    author_last_name: string;
}

export interface Event {
    id: number;
    title: string;
    description: string;
    eventDate: string;
    groupId: number;
}

export type GroupMemberReaction = {
    userID: number;
    fname: string;
    lname: string;
    reaction: 'going' | 'not going' | 'pending';
};
