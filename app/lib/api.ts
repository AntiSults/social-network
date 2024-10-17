import axios from 'axios';

// Base URL for the API
const API_URL = 'http://localhost:8080'; 

export const getPosts = async () => {
  try {
    // Send GET request to fetch posts
    const response = await axios.get(`${API_URL}/posts`, {
      withCredentials: true, // Include cookies with the request
    });
    return response.data;
  } catch (error) {
    console.error("Error fetching posts:", error);
    throw error;
  }
};

export const createPost = async (content: string, privacy: string, file: File | null, groupId?: number | null, visibleUsers?: number[] | null) => {
  try {
    const formData = new FormData();

    formData.append("content", content);
    formData.append("privacy", privacy);

    if (file) {
      formData.append("files", file);
    }
    if (groupId !== null) {
      formData.append("group_id", String(groupId));
  }
  if (visibleUsers && visibleUsers.length > 0) {
    formData.append("visible_users", JSON.stringify(visibleUsers));
  }

    // Send POST request to create a new post
    const response = await axios.post(`${API_URL}/create-posts`, formData, {
      withCredentials: true,
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  } catch (error) {
    console.error("Error creating post:", error);
    throw error;
  }
};

export const getComments = async (postId: number) => {
  try {
    const response = await axios.get(`${API_URL}/comments?post_id=${postId}`, {
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.error("Error fetching comments:", error);
    throw error;
  }
};

export const createComment = async (postId: number, content: string, file: File | null) => {
  try {
    const formData = new FormData();
    formData.append("post_id", postId.toString());
    formData.append("content", content);

    if (file) {
      formData.append("file", file);
    }

    const response = await axios.post(`${API_URL}/create-comment`, formData, {
      withCredentials: true,
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  } catch (error) {
    console.error("Error creating comment:", error);
    throw error;
  }
};

export const getUserGroups = async () => {
  try {
    const response = await axios.get(`${API_URL}/groups/get-users`, {
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.error("Error fetching user groups:", error);
    throw error;
  }
};