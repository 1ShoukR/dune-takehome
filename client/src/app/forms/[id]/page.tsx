'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import AuthenticatedLayout from '../../../components/AuthenticatedLayout';
import FormBuilder, { FormData } from '../../../components/FormBuilder/FormBuilder';
import { CreateFormRequest, formsAPI } from '../../../services/api';

export default function FormViewPage() {
  const params = useParams();
  const router = useRouter();
  const [form, setForm] = useState<FormData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    if (params.id) {
      fetchForm();
    }
  }, [params.id]);

  const fetchForm = async () => {
    try {
      setIsLoading(true);
      const response = await formsAPI.getFormById(params.id as string);
      setForm(response);
    } catch (error) {
      setError('Failed to load form');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSaveForm = async (formData: FormData, isDraft: boolean = true) => {
    try {
      setIsSaving(true);
      setError(null);
      
      const payload = {
        ...formData,
        status: isDraft ? 'draft' : 'published'
      };

      await formsAPI.updateForm(params.id as string, payload as CreateFormRequest);
      setIsEditing(false);
      fetchForm(); // Refresh form data
    } catch (error) {
      setError('Failed to save form. Please try again.');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <AuthenticatedLayout>
        <div className="flex items-center justify-center min-h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
        </div>
      </AuthenticatedLayout>
    );
  }

  if (error || !form) {
    return (
      <AuthenticatedLayout>
        <div className="text-center py-12">
          <p className="text-red-600">{error || 'Form not found'}</p>
          <button
            onClick={() => router.push('/forms')}
            className="mt-4 text-indigo-600 hover:text-indigo-900"
          >
            ← Back to Forms
          </button>
        </div>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <div className="max-w-6xl mx-auto">
        <div className="mb-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-black">
                {isEditing ? 'Edit Form' : 'View Form'}
              </h1>
              <p className="mt-1 text-sm text-black">
                {form.title}
              </p>
            </div>
            <div className="flex items-center space-x-3">
              {!isEditing && (
                <>
              <button
                onClick={() => router.push('/forms')}
                className="text-black hover:text-black"
              >
                ← Back to Forms
              </button>
                  <button
                    onClick={() => setIsEditing(true)}
                    className="px-4 py-2 text-sm font-medium text-black bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                  >
                    Edit Form
                  </button>
                  {form.status === 'published' && form.share_url && (
                    <button
                      onClick={() => {
                        navigator.clipboard.writeText(`${window.location.origin}/f/${form.share_url}`);
                        alert('Share link copied!');
                      }}
                      className="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700"
                    >
                      Copy Share Link
                    </button>
                  )}
                </>
              )}
            </div>
          </div>
        </div>

        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 rounded-md p-4">
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}

        {isEditing ? (
          <FormBuilder 
            onSave={handleSaveForm}
            isLoading={isSaving}
            initialData={form}
          />
        ) : (
          // Read-only view
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="max-w-2xl">
              <div className="mb-6">
                <h1 className="text-2xl font-bold text-black">{form.title}</h1>
                {form.description && (
                  <p className="mt-2 text-black">{form.description}</p>
                )}
              </div>
              
              <div className="space-y-6">
                {form.fields.map((field) => (
                  <div key={field.id} className="space-y-2">
                    <label className="block text-sm font-medium text-black">
                      {field.label}
                      {field.required && <span className="text-red-500 ml-1">*</span>}
                    </label>
                    
                    {field.type === 'text' || field.type === 'email' ? (
                      <input
                        type={field.type}
                        placeholder={field.placeholder}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 placeholder-gray-500 disabled:placeholder-gray-500"
                        disabled
                      />
                    ) : field.type === 'textarea' ? (
                      <textarea
                        placeholder={field.placeholder}
                        rows={3}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 placeholder-gray-500 disabled:placeholder-gray-500"
                        disabled
                      />
                    ) : field.type === 'rating' ? (
                      <div className="flex space-x-1">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <span key={star} className="text-2xl text-black">⭐</span>
                        ))}
                      </div>
                    ) : null}
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </AuthenticatedLayout>
  );
}
