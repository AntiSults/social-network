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

export const createPost = async (content: string, privacy: string) => {
  try {
    // Send POST request to create a new post
    const response = await axios.post(`${API_URL}/create-posts`, { content, privacy }, {
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    console.error("Error creating post:", error);
    throw error;
  }
};
