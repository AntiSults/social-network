import React from "react";

export interface Post {
    id: number;
    content: string;
    created_at: string;
    author_first_name: string;
    author_last_name: string;
}

interface PostsListProps {
    posts: Post[];
}

const PostsList: React.FC<PostsListProps> = ({ posts }) => {
    return (
        <div className="mt-10 flex flex-col items-center">
            {posts.map(post => (
                <div 
                  key={post.id} 
                  className="w-full max-w-lg p-6 bg-white shadow-md rounded-lg mb-6"
                >
                    <p className="text-lg font-semibold">
                        {post.author_first_name} {post.author_last_name}:
                    </p>
                    <p className="mt-2">{post.content}</p>
                    <small className="block mt-4 text-gray-500">
                        {new Date(post.created_at).toLocaleString()}
                    </small>
                </div>
            ))}
        </div>
    );
};

export default PostsList;
