import React, { useState, useEffect } from "react";
import { getComments, createComment } from "@/lib/api";
import { Comment } from "@/app/utils/types";

interface CommentsProps {
  postID: number;
}

const Comments: React.FC<CommentsProps> = ({ postID }) => {
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState<string>("");
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    const fetchComments = async () => {
      setLoading(true);
      try {
        const response = await getComments(postID);
        console.log("Fetched comments:", response);
        setComments(response);
      } catch (error) {
        console.error("Error fetching comments:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchComments();
  }, [postID]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (newComment.trim()) {
      setLoading(true);
      try {
        const response = await createComment(postID, newComment, file || null);
        if (response) {
          setNewComment("");
          setFile(null);
          const commentsResponse = await getComments(postID);
          setComments(commentsResponse);
        }
      } catch (error) {
        console.error("Error posting comment:", error);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div className="mt-4">
      <h3 className="text-lg font-bold mb-2">Comments</h3>
      {loading ? (
        <p>Loading comments...</p>
      ) : (
        <ul className="mb-4">
          {!comments || comments.length === 0 ? (
            <li className="text-gray-500">Be the first to comment!</li>
          ) : (
            comments.map((comment) => (
              <li key={comment.id} className="border-b p-2">
                <p className="text-sm">
                  <strong>{comment.author_first_name} {comment.author_last_name}:</strong> {comment.content}
                </p>
                {comment.file && (
                  <div className="mt-4">
                    <img
                      src={comment.file.replace("../public", "")}
                      alt="Comment file"
                      className="mt-2 max-w-full rounded"
                      onError={(e) => {
                        const target = e.target as HTMLImageElement;
                        target.onerror = null;
                      }}
                    />
                  </div>
                )}
                <small className="block mt-4 text-gray-500">
                  {new Date(comment.created_at).toLocaleString()}
                </small>
              </li>
            ))
          )}
        </ul>
      )}

      <form onSubmit={handleSubmit} className="relative mx-auto max-w-lg p-4 bg-white shadow-md rounded-lg mb-4">
        <div className="mb-4">
          <input
            type="text"
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Add a comment..."
            className="w-full p-2 border rounded-md border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="mb-4">
          <input
            type="file"
            accept="image/*, .gif"
            onChange={(e) => setFile(e.target.files ? e.target.files[0] : null)}
            className="border rounded-md border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 w-full"
          />
          <button
            type="submit"
            className="bg-gray-600 text-white px-4 py-2 rounded hover:bg-gray-700 w-full"
            disabled={loading}
          >
            {loading ? "Submitting..." : "Comment"}
          </button>
        </div>
      </form>
    </div>
  );
  
  
  
};

export default Comments;
