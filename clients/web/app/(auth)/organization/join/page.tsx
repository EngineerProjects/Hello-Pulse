'use client';

import React from 'react';
import JoinOrgForm from '@/components/auth/JoinOrgForm';

export default function JoinOrganizationPage() {
  return (
    <div className="h-full bg-gray-100 dark:bg-gray-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-[700px]">
        <h2 className="mt-6 text-center text-5xl font-extrabold text-gray-900 dark:text-white">
          Join or create an organization
        </h2>
        <p className="mt-2 text-center text-gray-600 max-w text-xl dark:text-gray-400">
          To use Hello Pulse, you need to join an organization. <br />
          If you have an invitation code, you can enter it below. <br />
          This code can be obtained from your organization's administrator.
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <JoinOrgForm />
      </div>
    </div>
  );
}