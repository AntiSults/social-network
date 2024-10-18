import React, { useEffect, useState } from 'react';
import { User, Post } from '@/app/utils/types';
import Image from 'next/image';
import FollowList from './FollowLists';
import PostsList from './PostsList';
import { getUserPosts } from '@/app/lib/api';


interface Props {
    foundUser: User;
    currentUser: User | null;
}

const FoundUserProfile: React.FC<Props> = ({ foundUser, currentUser }) => {
    const [isFollower, setIsFollower] = useState<boolean>(false);

    const [canViewFullProfile, setCanViewFullProfile] = useState<boolean>(false);

    const [userPosts, setUserPosts] = useState<Post[]>([]);


    useEffect(() => {
        const checkFollowerStatus = async () => {
            if (currentUser && foundUser) {
                try {
                    const res = await fetch(`http://localhost:8080/followers/status?userId=${foundUser.ID}&followerId=${currentUser.ID}`);
                    if (!res.ok) {
                        throw new Error(`Error: ${res.status}`);
                    }
                    const data = await res.json();
                    setIsFollower(data.isFollowing);
                } catch (error) {
                    console.error('Error fetching follower status:', error);
                }
            }
        };
        if (currentUser && foundUser) {
            checkFollowerStatus();
        }
    }, [currentUser, foundUser]);

    useEffect(() => {
        const fetchUserPosts = async () => {
          try {
            const postsData = await getUserPosts(foundUser.ID);
            setUserPosts(postsData);
          } catch (error) {
            console.error('Error fetching user posts:', error);
          }
        };
      
        if (foundUser) {
          fetchUserPosts();
        }
      }, [foundUser]);

    useEffect(() => {
        if (foundUser.profileVisibility === "public" || (foundUser.profileVisibility === "private" && isFollower)) {
            setCanViewFullProfile(true);
        } else {
            setCanViewFullProfile(false);
        }
    }, [foundUser, isFollower]);

    return (

        <div className="bg-white shadow-md rounded-lg p-8 max-w-lg w-full text-center">
            <h1 className="text-2xl font-bold mb-4">
                {`${foundUser.firstName} ${foundUser.lastName}'s Profile`}
            </h1>
            {canViewFullProfile ? (
                <div className="flex flex-col items-center">
                    <Image
                        src={foundUser.avatarPath ? `${foundUser.avatarPath}` : "/default_avatar.jpg"}
                        alt={`${foundUser.firstName}'s Avatar`}
                        width={250}
                        height={250}
                        className="rounded-full shadow-lg"
                    />
                    <p className="text-gray-600 mt-4">
                        Made in: {foundUser.dob ? new Date(foundUser.dob).toLocaleDateString('en-US', {
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric',
                        }) : "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        Nick: {foundUser.nickname || "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        About me: {foundUser.aboutMe || "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        Profile: {foundUser.profileVisibility || "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        Email: {foundUser.email || "No details provided"}
                    </p>
                    {<FollowList user={foundUser} />}

                    <div className="mt-8 w-full">
                    <h2 className="text-xl font-bold mb-4">Posts by {foundUser.firstName}</h2>
                    <PostsList posts={userPosts} />
                </div>
                </div>
            ) : (
                <p>This profile is private. You can only see limited information.</p>
            )}
        </div>
    );
};

export default FoundUserProfile;
