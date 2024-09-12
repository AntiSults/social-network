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

export const createPost = async (content: string, privacy: string, file: File | null) => {
  try {
    const formData = new FormData();
    formData.append("content", content);
    formData.append("privacy", privacy);

    if (file) {
      formData.append("files", file);
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
