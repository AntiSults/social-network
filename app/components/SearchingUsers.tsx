import { useState } from 'react';
import { searchUsers } from '@/app/utils/searchUsers';
import { User as SearchResult } from '@/app/utils/types';
import FieldInput from './FieldInput';
import Button from './Button';
import { useRouter } from 'next/navigation'; // Import useRouter from next/navigation

const UserSearch: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
    const router = useRouter();  // Initialize router

    const handleSearch = async (e: React.FormEvent) => {
        e.preventDefault();
        const users = await searchUsers(searchQuery);
        setSearchResults(users || []);
    };

    const handleUserClick = (user: SearchResult) => {
        // Navigate to the found user's profile page using their ID
        router.push(`/users/${user.ID}/profile`);  // Change the URL to /profile/userID
    };

    return (
        <div>
            <form onSubmit={handleSearch}>
                <FieldInput
                    type='text'
                    placeholder='Search for users...'
                    required={true}
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                />
                <Button type='submit' name='Search' />
            </form>
            <ul>
                {searchResults.map(result => (
                    <li key={result.ID} onClick={() => handleUserClick(result)} className='cursor-pointer hover:bg-gray-100 p-2'>
                        <p className='font-semibold'>
                            {result.firstName} {result.lastName}
                        </p>
                        <p className='text-sm text-gray-500'>{result.email}</p>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default UserSearch;
