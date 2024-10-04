import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { UserProvider } from "./context/UserContext"; // Import the UserProvider

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "social-network",
  description: "Social Network project for Kood/Jõhvi",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        {/* Wrap children with UserProvider for context */}
        <UserProvider>
          {children}
        </UserProvider>
      </body>
    </html>
  );
}