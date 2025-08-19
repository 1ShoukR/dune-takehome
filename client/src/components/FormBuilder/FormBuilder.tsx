'use client';

import { useState, useCallback } from 'react';

export interface FormField {
  id: string;
  type: 'text' | 'textarea' | 'email' | 'number' | 'select' | 'radio' | 'checkbox' | 'rating';
  label: string;
  placeholder?: string;
  required: boolean;
  options?: string[];
  validation?: {
    minLength?: number;
    maxLength?: number;
    min?: number;
    max?: number;
    pattern?: string;
  };
  order: number;
}

export interface FormData {
  title: string;
  description: string;
  fields: FormField[];
}

interface FormBuilderProps {
  onSave: (formData: FormData, isDraft?: boolean) => void;
  isLoading: boolean;
  initialData?: FormData;
}

export default function FormBuilder({ onSave, isLoading, initialData }: FormBuilderProps) {
  const [formData, setFormData] = useState<FormData>(
    initialData || {
      title: '',
      description: '',
      fields: [],
    }
  );
  const [selectedField, setSelectedField] = useState<string | null>(null);
  const [previewMode, setPreviewMode] = useState(false);

  const updateFormSettings = useCallback((updates: Partial<FormData>) => {
    setFormData(prev => ({ ...prev, ...updates }));
  }, []);

  const addField = useCallback((fieldType: FormField['type']) => {
    const newField: FormField = {
      id: `field_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      type: fieldType,
      label: `New ${fieldType.charAt(0).toUpperCase() + fieldType.slice(1)} Field`,
      placeholder: fieldType === 'text' || fieldType === 'email' || fieldType === 'textarea' 
        ? 'Enter your answer...' 
        : undefined,
      required: false,
      options: ['select', 'radio', 'checkbox'].includes(fieldType) 
        ? ['Option 1', 'Option 2'] 
        : undefined,
      order: formData.fields.length,
    };

    setFormData(prev => ({
      ...prev,
      fields: [...prev.fields, newField],
    }));
    setSelectedField(newField.id);
  }, [formData.fields.length]);

  const updateField = useCallback((fieldId: string, updates: Partial<FormField>) => {
    setFormData(prev => ({
      ...prev,
      fields: prev.fields.map(field =>
        field.id === fieldId ? { ...field, ...updates } : field
      ),
    }));
  }, []);

  const deleteField = useCallback((fieldId: string) => {
    setFormData(prev => ({
      ...prev,
      fields: prev.fields.filter(field => field.id !== fieldId)
        .map((field, index) => ({ ...field, order: index })),
    }));
    setSelectedField(null);
  }, []);

  const moveField = useCallback((fieldId: string, direction: 'up' | 'down') => {
    setFormData(prev => {
      const currentIndex = prev.fields.findIndex(f => f.id === fieldId);
      if (currentIndex === -1) return prev;
      
      const newIndex = direction === 'up' ? currentIndex - 1 : currentIndex + 1;
      if (newIndex < 0 || newIndex >= prev.fields.length) return prev;
      
      const newFields = [...prev.fields];
      [newFields[currentIndex], newFields[newIndex]] = [newFields[newIndex], newFields[currentIndex]];
      
      return {
        ...prev,
        fields: newFields.map((field, index) => ({ ...field, order: index })),
      };
    });
  }, []);

  const handleSave = (isDraft: boolean = true) => {
    if (!formData.title.trim()) {
      alert('Please enter a form title');
      return;
    }
    onSave(formData, isDraft);
  };

  const canSave = formData.title.trim() && formData.fields.length > 0;

  // Field type options for the palette
  const fieldTypes = [
    { type: 'text', label: 'Text Input', icon: '📝' },
    { type: 'textarea', label: 'Long Text', icon: '📄' },
    { type: 'email', label: 'Email', icon: '📧' },
    { type: 'number', label: 'Number', icon: '🔢' },
    { type: 'select', label: 'Dropdown', icon: '📋' },
    { type: 'radio', label: 'Multiple Choice', icon: '⚪' },
    { type: 'checkbox', label: 'Checkboxes', icon: '☑️' },
    { type: 'rating', label: 'Rating', icon: '⭐' },
  ] as const;

  return (
    <div className="grid grid-cols-12 gap-6 min-h-screen">
      {/* Left Sidebar - Field Palette & Settings */}
      <div className="col-span-3 space-y-6">
        {/* Field Palette */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Add Fields</h3>
          <div className="space-y-2">
            {fieldTypes.map(({ type, label, icon }) => (
              <button
                key={type}
                onClick={() => addField(type)}
                className="w-full flex items-center p-3 text-left border border-gray-200 rounded-lg hover:bg-gray-50 hover:border-indigo-300 transition-colors"
              >
                <span className="text-lg mr-3">{icon}</span>
                <span className="text-sm font-medium text-gray-700">{label}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Form Settings */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Form Settings</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Form Title *
              </label>
              <input
                type="text"
                value={formData.title}
                onChange={(e) => updateFormSettings({ title: e.target.value })}
                placeholder="Enter form title..."
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                value={formData.description}
                onChange={(e) => updateFormSettings({ description: e.target.value })}
                placeholder="Enter form description..."
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Main Canvas Area */}
      <div className="col-span-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="flex items-center justify-between p-4 border-b border-gray-200">
            <div className="flex items-center space-x-2">
              <button
                onClick={() => setPreviewMode(false)}
                className={`px-3 py-1 text-sm font-medium rounded ${
                  !previewMode
                    ? 'bg-indigo-100 text-indigo-700'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                Build
              </button>
              <button
                onClick={() => setPreviewMode(true)}
                className={`px-3 py-1 text-sm font-medium rounded ${
                  previewMode
                    ? 'bg-indigo-100 text-indigo-700'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                Preview
              </button>
            </div>
            
            <div className="flex items-center space-x-2">
              <button
                onClick={() => handleSave(true)}
                disabled={!canSave || isLoading}
                className="px-3 py-1 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
              >
                {isLoading ? 'Saving...' : 'Save Draft'}
              </button>
              <button
                onClick={() => handleSave(false)}
                disabled={!canSave || isLoading}
                className="px-3 py-1 text-sm font-medium text-white bg-indigo-600 rounded hover:bg-indigo-700 disabled:opacity-50"
              >
                Publish
              </button>
            </div>
          </div>

          <div className="p-6">
            {previewMode ? (
              // Preview Mode
              <div className="max-w-2xl">
                <div className="mb-6">
                  <h1 className="text-2xl font-bold text-gray-900">
                    {formData.title || 'Untitled Form'}
                  </h1>
                  {formData.description && (
                    <p className="mt-2 text-gray-600">{formData.description}</p>
                  )}
                </div>
                
                <div className="space-y-6">
                  {formData.fields.length === 0 ? (
                    <p className="text-gray-500 text-center py-8">
                      No fields added yet. Switch to Build mode to add fields.
                    </p>
                  ) : (
                    formData.fields.map((field) => (
                      <div key={field.id} className="space-y-2">
                        <label className="block text-sm font-medium text-gray-700">
                          {field.label}
                          {field.required && <span className="text-red-500 ml-1">*</span>}
                        </label>
                        
                        {field.type === 'text' || field.type === 'email' ? (
                          <input
                            type={field.type}
                            placeholder={field.placeholder}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                            disabled
                          />
                        ) : field.type === 'textarea' ? (
                          <textarea
                            placeholder={field.placeholder}
                            rows={3}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                            disabled
                          />
                        ) : field.type === 'number' ? (
                          <input
                            type="number"
                            placeholder={field.placeholder}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                            disabled
                          />
                        ) : field.type === 'select' ? (
                          <select className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" disabled>
                            <option>Select an option...</option>
                            {field.options?.map((option, index) => (
                              <option key={index} value={option}>{option}</option>
                            ))}
                          </select>
                        ) : field.type === 'radio' ? (
                          <div className="space-y-2">
                            {field.options?.map((option, index) => (
                              <div key={index} className="flex items-center">
                                <input type="radio" name={field.id} className="mr-2" disabled />
                                <span className="text-sm text-gray-700">{option}</span>
                              </div>
                            ))}
                          </div>
                        ) : field.type === 'checkbox' ? (
                          <div className="space-y-2">
                            {field.options?.map((option, index) => (
                              <div key={index} className="flex items-center">
                                <input type="checkbox" className="mr-2" disabled />
                                <span className="text-sm text-gray-700">{option}</span>
                              </div>
                            ))}
                          </div>
                        ) : field.type === 'rating' ? (
                          <div className="flex space-x-1">
                            {[1, 2, 3, 4, 5].map((star) => (
                              <span key={star} className="text-2xl text-gray-300 cursor-pointer">⭐</span>
                            ))}
                          </div>
                        ) : null}
                      </div>
                    ))
                  )}
                </div>
              </div>
            ) : (
              // Build Mode
              <div>
                <div className="mb-6">
                  <h1 className="text-2xl font-bold text-gray-900">
                    {formData.title || 'Untitled Form'}
                  </h1>
                  {formData.description && (
                    <p className="mt-2 text-gray-600">{formData.description}</p>
                  )}
                </div>
                
                <div className="space-y-4">
                  {formData.fields.length === 0 ? (
                    <div className="text-center py-12 border-2 border-dashed border-gray-300 rounded-lg">
                      <p className="text-gray-500">
                        Add fields from the left panel to start building your form
                      </p>
                    </div>
                  ) : (
                    formData.fields.map((field, index) => (
                      <div
                        key={field.id}
                        onClick={() => setSelectedField(field.id)}
                        className={`p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                          selectedField === field.id
                            ? 'border-indigo-500 bg-indigo-50'
                            : 'border-gray-200 hover:border-gray-300'
                        }`}
                      >
                        <div className="flex items-center justify-between">
                          <div className="flex-1">
                            <div className="flex items-center space-x-2 mb-2">
                              <span className="text-sm font-medium text-gray-900">
                                {field.label}
                              </span>
                              {field.required && (
                                <span className="text-red-500 text-sm">*</span>
                              )}
                              <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
                                {field.type}
                              </span>
                            </div>
                            
                            {/* Field Preview */}
                            <div className="text-sm text-gray-600">
                              {field.placeholder && (
                                <span className="italic">Placeholder: {field.placeholder}</span>
                              )}
                              {field.options && (
                                <span>Options: {field.options.join(', ')}</span>
                              )}
                            </div>
                          </div>
                          
                          <div className="flex items-center space-x-2">
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                moveField(field.id, 'up');
                              }}
                              disabled={index === 0}
                              className="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-50"
                            >
                              ↑
                            </button>
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                moveField(field.id, 'down');
                              }}
                              disabled={index === formData.fields.length - 1}
                              className="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-50"
                            >
                              ↓
                            </button>
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                deleteField(field.id);
                              }}
                              className="p-1 text-red-400 hover:text-red-600"
                            >
                              🗑️
                            </button>
                          </div>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Right Sidebar - Field Properties */}
      <div className="col-span-3">
        {selectedField && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              Field Properties
            </h3>
            
            {(() => {
              const field = formData.fields.find(f => f.id === selectedField);
              if (!field) return null;
              
              return (
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Label
                    </label>
                    <input
                      type="text"
                      value={field.label}
                      onChange={(e) => updateField(field.id, { label: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                    />
                  </div>
                  
                  {(field.type === 'text' || field.type === 'email' || field.type === 'textarea' || field.type === 'number') && (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Placeholder
                      </label>
                      <input
                        type="text"
                        value={field.placeholder || ''}
                        onChange={(e) => updateField(field.id, { placeholder: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                      />
                    </div>
                  )}
                  
                  <div className="flex items-center">
                    <input
                      type="checkbox"
                      checked={field.required}
                      onChange={(e) => updateField(field.id, { required: e.target.checked })}
                      className="mr-2"
                    />
                    <label className="text-sm font-medium text-gray-700">
                      Required field
                    </label>
                  </div>
                  
                  {(field.type === 'select' || field.type === 'radio' || field.type === 'checkbox') && (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Options (one per line)
                      </label>
                      <textarea
                        value={field.options?.join('\n') || ''}
                        onChange={(e) => updateField(field.id, { 
                          options: e.target.value.split('\n').filter(opt => opt.trim()) 
                        })}
                        rows={4}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                      />
                    </div>
                  )}
                </div>
              );
            })()}
          </div>
        )}
      </div>
    </div>
  );
}