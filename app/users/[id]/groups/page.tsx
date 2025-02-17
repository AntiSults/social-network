"use client";
import GroupList from '@/app/components/GroupList';
import NavBar from '@/app/components/NavBar';
import CreateGroupForm from '@/app/components/CreateGroupForm';
import UserSearch from '@/app/components/SearchingUsers';
import InviteToGroup from '@/app/components/InviteToGroup';
import { useState } from 'react';
import { useUser } from '@/app/context/UserContext';
import { User } from '@/app/utils/types';

const GroupsPage = () => {
    const { user: currentUser } = useUser();

    const [selectedUser, setSelectedUser] = useState<User | null>(null);
    const [selectedGroup, setSelectedGroup] = useState<number | null>(null);

    if (!currentUser) {
        return (
            <div className="min-h-screen">
                <NavBar logged={false} />
                <p className="text-center text-gray-600">Please login to see Groups!</p>
            </div>
        );
    }

    return (
        <div className="min-h-screen">
            <NavBar logged={true} />

            <div className="container mx-auto p-6">
                <h1 className="text-3xl font-bold mb-6 text-center">Groups</h1>

                <div className="bg-white shadow-md rounded-lg p-6 mb-6">
                    <h2 className="text-2xl font-semibold mb-4">Create New Group</h2>
                    <CreateGroupForm />
                </div>

                <div className="bg-white shadow-md rounded-lg p-6 mb-6">
                    <h2 className="text-2xl font-semibold mb-4">Available Groups</h2>
                    <GroupList onSelectGroup={setSelectedGroup} actionType="invite" />

                </div>

                <div className="bg-white shadow-md rounded-lg p-6 mb-6">
                    <h2 className="text-2xl font-semibold mb-4">Invite User to Group</h2>
                    <UserSearch onSelectUser={setSelectedUser} />

                    {selectedUser && selectedGroup && (
                        <InviteToGroup
                            groupId={selectedGroup}
                            invitedUser={selectedUser}
                            currentUser={currentUser}
                        />
                    )}
                </div>
            </div>
        </div>
    );
};

export default GroupsPage;
