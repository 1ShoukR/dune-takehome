'use client';

import { useState, useEffect } from 'react';
import AuthenticatedLayout from '../../components/AuthenticatedLayout';
import { formsAPI } from '../../services/api';
import { useRouter } from 'next/navigation';

interface Form {
  id: string;
  title: string;
  description?: string;
  status: 'draft' | 'published' | 'archived';
  share_url?: string; 
  created_at: string;
  updated_at: string;
}

export default function FormsPage() {
  const [forms, setForms] = useState<Form[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    fetchForms();
  }, []);

  const fetchForms = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await formsAPI.getAllForms();
      setForms(response.forms || []);
    } catch (error: unknown) {
      console.error('Error fetching forms:', error);
      setError('Failed to load forms. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleNewForm = () => {
    router.push('/forms/new');
  };

  const handleViewDrafts = () => {
    // TODO: Filter to show only drafts
    console.log('Filter to show drafts only');
  };

  return (
    <AuthenticatedLayout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Your Forms</h1>
            <p className="mt-1 text-sm text-gray-600">
              Create and manage your custom forms
            </p>
          </div>
          
          <div className="flex space-x-3">
            <button
              onClick={handleViewDrafts}
              className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <svg className="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              View Drafts
            </button>
            <button
              onClick={handleNewForm}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <svg className="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              New Form
            </button>
          </div>
        </div>

        {isLoading && (
          <div className="bg-white shadow rounded-lg p-6">
            <div className="flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
              <span className="ml-3 text-gray-600">Loading forms...</span>
            </div>
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-sm text-red-800">{error}</p>
              </div>
            </div>
          </div>
        )}

        {!isLoading && !error && (
          <div className="bg-white shadow rounded-lg">
            {forms.length === 0 ? (
              <div className="text-center py-12">
                <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <h3 className="mt-2 text-sm font-medium text-gray-900">No saved forms</h3>
                <p className="mt-1 text-sm text-gray-500">
                  Get started by creating your first form.
                </p>
                <div className="mt-6">
                  <button
                    onClick={handleNewForm}
                    className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                  >
                    <svg className="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    Create New Form
                  </button>
                </div>
              </div>
            ) : (
              <div className="divide-y divide-gray-200">
                {forms.map((form) => (
                  <div
                    key={form.id}
                    className="p-6 hover:bg-gray-50 cursor-pointer"
                    onClick={() => router.push(`/forms/${form.id}`)}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <h3 className="text-lg font-medium text-gray-900">
                          {form.title}
                        </h3>
                        {form.description && (
                          <p className="mt-1 text-sm text-gray-500">
                            {form.description}
                          </p>
                        )}
                        <div className="mt-2 flex items-center space-x-4 text-sm text-gray-500">
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            form.status === 'published' 
                              ? 'bg-green-100 text-green-800'
                              : form.status === 'draft'
                              ? 'bg-yellow-100 text-yellow-800'
                              : 'bg-gray-100 text-gray-800'
                          }`}>
                            {form.status.charAt(0).toUpperCase() + form.status.slice(1)}
                          </span>
                          <span>
                            Updated {new Date(form.updated_at).toLocaleDateString()}
                          </span>
                          {form.status === 'published' && form.share_url && (
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                navigator.clipboard.writeText(`${window.location.origin}/f/${form.share_url}`);
                                alert('Share link copied to clipboard!');
                              }}
                              className="text-green-600 hover:text-green-900 text-sm font-medium"
                            >
                              Share
                            </button>
                          )}
                        </div>
                      </div>
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            router.push(`/forms/${form.id}/edit`);
                          }}
                          className="text-indigo-600 hover:text-indigo-900 text-sm font-medium"
                        >
                          Edit
                        </button>
                        {form.status === 'published' && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              navigator.clipboard.writeText(`${window.location.origin}/f/${form.share_url}`);
                              alert('Share link copied!');
                            }}
                            className="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700"
                          >
                            Copy Share Link
                          </button>
                        )}
                        <button className="text-gray-400 hover:text-gray-600 text-sm font-medium">
                          •••
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </AuthenticatedLayout>
  );
}