"use client";

import React, { useEffect, useState, useCallback } from "react";
import { getPosts, createPost } from "@/app/lib/api";
import NewPostForm from "@/app/components/NewPostForm";
import PostsList from "@/app/components/PostsList";
import checkLoginStatus from "@/app/utils/checkLoginStatus";
import NavBar from "@/app/components/NavBar";
import Comments from "@/app/components/Comments";

const PostsPage = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [posts, setPosts] = useState<any[]>([]);

  const fetchPosts = useCallback(async () => {
    try {
      const fetchedPosts = await getPosts();
      setPosts(fetchedPosts && Array.isArray(fetchedPosts) ? fetchedPosts : []);
    } catch (error) {
      console.error("Failed to fetch posts:", error);
    }
  }, []);

  useEffect(() => {
    setIsLoggedIn(checkLoginStatus());
    fetchPosts();
  }, [fetchPosts]);

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
        {posts.map((post) => (
          <div key={post.id}>
            <h3>{post.content}</h3>
            <Comments postID={post.id} />
          </div>
        ))}
      </div>
    </>
  );
};

export default PostsPage;
