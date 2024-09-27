"use client";
import GroupList from '@/app/components/GroupList';
import NavBar from "@/app/components/NavBar";
import CreateGroupForm from '@/app/components/CreateGroupForm';
import UserSearch from '@/app/components/SearchingUsers';
import InviteToGroup from '@/app/components/InviteToGroup'; // Import the InviteToGroup component
import { useState } from 'react';
import { useUser } from '@/app/context/UserContext';
import { Group, User } from '@/app/utils/types';

const GroupsPage = () => {
    const { user: currentUser } = useUser();
    const [selectedUser, setSelectedUser] = useState<User | null>(null);
    const [selectedGroup, setSelectedGroup] = useState<number | null>(null);

    return (
        <div className="min-h-screen bg-gray-50">
            <NavBar logged={true} />

            <div>
                <h1>Groups</h1>
                {/* 1. Display all groups with the option to join */}
                <GroupList onSelectGroup={setSelectedGroup} /> {/* Pass selected group ID */}

                {/* 2. Create a new group */}
                <div>
                    <h2>Create New Group</h2>
                    <CreateGroupForm />
                </div>

                {/* 3. Invite users to groups */}
                <div>
                    <h2>Invite User to Group</h2>
                    <UserSearch onSelectUser={setSelectedUser} /> {/* Capture selected user */}

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
