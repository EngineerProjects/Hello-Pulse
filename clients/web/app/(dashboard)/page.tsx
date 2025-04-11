import React from 'react';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Dashboard | Hello Pulse',
  description: 'Hello Pulse dashboard',
};

export default function DashboardPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        Welcome to Hello Pulse
      </h1>
      <p className="text-gray-600 dark:text-gray-400">
        Your collaborative brainstorming platform powered by AI.
      </p>
      
      {/* Placeholder for dashboard content */}
      <div className="mt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h2 className="text-lg font-semibold mb-4">Recent Projects</h2>
          <p className="text-gray-500 dark:text-gray-400">No projects yet</p>
        </div>
        
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h2 className="text-lg font-semibold mb-4">Upcoming Events</h2>
          <p className="text-gray-500 dark:text-gray-400">No events scheduled</p>
        </div>
        
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
          <h2 className="text-lg font-semibold mb-4">Team Members</h2>
          <p className="text-gray-500 dark:text-gray-400">No team members</p>
        </div>
      </div>
    </div>
  );
}