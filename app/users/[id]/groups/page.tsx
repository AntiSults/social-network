"use client"
import GroupList from '../../../components/GroupList';
import NavBar from "../../../components/NavBar";
import CreateGroupForm from '../../../components/CreateGroupForm';

import SearchBar from '../../../components/SearchBar';
import { useUser } from '../../../context/UserContext';

const GroupsPage = () => {
    const { user } = useUser();


    return (
        <div className="min-h-screen bg-gray-50">
            <NavBar logged={true} />

            <div>
                <h1>Groups</h1>
                {/* 1. Display all groups with the option to join */}
                <GroupList />
                {/* 2. Create a new group */}
                <div>
                    <h2>Create New Group</h2>
                    <CreateGroupForm />
                </div>
                {/* 3. Invite users to groups */}
                <div>
                    <h2>Invite Users to Group</h2>
                    <SearchBar />
                </div>
            </div>
        </div>
    );
};

export default GroupsPage;
