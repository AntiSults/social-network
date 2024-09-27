"use client";

import React, { useEffect, useState } from "react";
import { getPosts, createPost } from "@/lib/api";
import NewPostForm from "@/app/components/NewPostForm";
import PostsList from "@/app/components/PostsList";
import checkLoginStatus from "@/app/utils/checkLoginStatus";
import NavBar from "@/app/components/NavBar";

const PostsPage = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [posts, setPosts] = useState<any[]>([]);

  useEffect(() => {
    setIsLoggedIn(checkLoginStatus());
    fetchPosts();
  }, []);

  const fetchPosts = async () => {
    try {
      const fetchedPosts = await getPosts();
      setPosts(fetchedPosts);
    } catch (error) {
      console.error("Failed to fetch posts:", error);
    }
  };


  const handlePostCreated = async (content: string, privacy: string, file?: File | null) => {
    try {
      await createPost(content, privacy, file || null);
      fetchPosts();
    } catch (error) {
      console.error("Failed to create post:", error);
    }
  };

  return (
    <>
      <NavBar logged={isLoggedIn} />
      <div>
        {isLoggedIn ? (
          <NewPostForm onPostCreated={handlePostCreated} />
        ) : (
          <p className="text-center text-gray-600">Please log in to create a post.</p>
        )}
        <PostsList posts={posts} />
      </div>
    </>
  );
}

export default PostsPage;
