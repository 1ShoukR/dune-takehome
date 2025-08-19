'use client';

import AuthenticatedLayout from '../../components/AuthenticatedLayout';

export default function FormsPage() {
  return (
    <AuthenticatedLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Your Forms</h1>
          <p className="mt-1 text-sm text-gray-600">
            Create and manage your custom forms
          </p>
        </div>

        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-900 mb-4">
            Form Builder Coming Soon
          </h2>
          <p className="text-gray-600">
            This is where you&apos;ll be able to create and manage your forms.
          </p>
        </div>
      </div>
    </AuthenticatedLayout>
  );
}