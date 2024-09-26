import React from "react";
import { Post } from "../utils/types";

interface PostsListProps {
    posts: Post[];
}

const PostsList: React.FC<PostsListProps> = ({ posts }) => {
    return (
        <div className="mt-10 flex flex-col items-center">
            {!posts || posts.length === 0 ? (
                <p className="text-center text-gray-600"> Empty database </p>
            ) : (
                posts.map(post => (
                    <div
                        key={post.id}
                        className="w-full max-w-lg p-6 bg-white shadow-md rounded-lg mb-6"
                    >
                        <p className="text-lg font-semibold">
                            {post.author_first_name} {post.author_last_name}:
                        </p>
                        <p className="mt-2">{post.content}</p>
                        {post.files && (
                            <div className="mt-4">
                                <img
                                    src={post.files.replace("../public", "")}
                                    alt="Attachment"
                                    className="mt-2 max-w-full rounded"
                                    onError={(e) => {
                                        const target = e.target as HTMLImageElement;
                                    }}
                                />

                            </div>
                        )}
                        <small className="block mt-4 text-gray-500">
                            {new Date(post.created_at).toLocaleString()}
                        </small>
                    </div>
                ))
            )}
        </div>
    );
};

export default PostsList;
