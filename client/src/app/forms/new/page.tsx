'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import AuthenticatedLayout from '../../../components/AuthenticatedLayout';
import FormBuilder, {FormData} from '@/components/FormBuilder/FormBuilder';
import { CreateFormRequest, formsAPI } from '../../../services/api';

export default function NewFormPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSaveForm = async (formData: FormData, isDraft: boolean = true) => {
    try {
      setIsLoading(true);
      setError(null);
      
      const payload = {
        ...formData,
        status: isDraft ? 'draft' : 'published'
      };

      await formsAPI.createForm(payload as CreateFormRequest);
      router.push('/forms');
    } catch (error: unknown) {
      setError('Failed to save form. Please try again.');
      console.error('Error saving form:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthenticatedLayout>
      <div className="max-w-6xl mx-auto">
        <div className="mb-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Create New Form</h1>
              <p className="mt-1 text-sm text-gray-600">
                Build your custom form with drag-and-drop fields
              </p>
            </div>
            <button
              onClick={() => router.push('/forms')}
              className="text-gray-600 hover:text-gray-900"
            >
              ← Back to Forms
            </button>
          </div>
        </div>

        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 rounded-md p-4">
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}

        <FormBuilder 
          onSave={handleSaveForm}
          isLoading={isLoading}
        />
      </div>
    </AuthenticatedLayout>
  );
}