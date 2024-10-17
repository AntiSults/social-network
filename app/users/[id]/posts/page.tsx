"use client";

import React, { useEffect, useState } from "react";
import { getPosts, createPost } from "@/app/lib/api"; 
import NewPostForm from "@/app/components/NewPostForm";
import PostsList from "@/app/components/PostsList";
import checkLoginStatus from "@/app/utils/checkLoginStatus";
import NavBar from "@/app/components/NavBar";
import { useUser } from '@/app/context/UserContext';

const PostsPage = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [posts, setPosts] = useState<any[]>([]);
    
    const { user, selectedUser } = useUser();
    const profileUser = selectedUser || user;

    useEffect(() => {
        const fetchData = async () => {
            setIsLoggedIn(checkLoginStatus());
            const fetchedPosts = await getPosts();
            setPosts(fetchedPosts);
        };


        fetchData();
    }, []); 

    const handlePostCreated = async (content: string, privacy: string, file?: File | null, groupId?: number | null, visibleUsers?: number[] | null) => {

        try {
            await createPost(content, privacy, file || null, groupId, visibleUsers);
            const fetchedPosts = await getPosts();
            setPosts(fetchedPosts);
        } catch (error) {
            console.error("Failed to create post:", error);
        }
    };
    

    return (
        <>
            <NavBar logged={isLoggedIn} />
            <div>
                {isLoggedIn ? (
                    <NewPostForm onPostCreated={handlePostCreated} user={user} />
                ) : (
                    <p className="text-center text-gray-600">Please log in to create a post.</p>
                )}
                <PostsList posts={posts} />
            </div>
        </>
    );
};

export default PostsPage;
