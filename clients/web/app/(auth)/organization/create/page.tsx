'use client';

import React from 'react';
import CreateOrgForm from '@/components/auth/CreateOrgForm';

export default function CreateOrganizationPage() {
  return (
    <div className="h-full bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8 dark:bg-gray-900">
      <div className="sm:mx-auto sm:w-full sm:max-w-[700px]">
        <h2 className="mt-6 text-center text-5xl font-extrabold text-gray-900 dark:text-white">
          Create an organization
        </h2>
        <p className="mt-2 text-center text-gray-600 max-w text-xl dark:text-gray-400">
          To proceed please enter the name of your organization.
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <CreateOrgForm />
      </div>
    </div>
  );
}