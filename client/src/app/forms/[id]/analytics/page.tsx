/* eslint-disable @typescript-eslint/no-explicit-any */
'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import AuthenticatedLayout from '../../../../components/AuthenticatedLayout';
import { useSocket } from '../../../../hooks/useSocket';
import api from '@/services/api';

interface FieldAnalytics {
  field_id: string;
  field_label: string;
  field_type: string;
  response_count: number;
  data: any;
}

interface FormAnalytics {
  form_id: string;
  form_title: string;
  total_responses: number;
  field_analytics: FieldAnalytics[];
  created_at: string;
}

export default function AnalyticsPage() {
  const params = useParams();
  const router = useRouter();
  const [analytics, setAnalytics] = useState<FormAnalytics | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());
  const [hasJoinedRoom, setHasJoinedRoom] = useState(false);
  
  // WebSocket connection
  const socket = useSocket('/');

  useEffect(() => {
    fetchAnalytics();
  }, [params.id]);

  useEffect(() => {
    if (socket && params.id && !hasJoinedRoom) {
      const tryJoinRoom = () => {
        socket.emit('join-analytics', { form_id: params.id });
        console.log('📊 Attempting to join analytics room for form:', params.id);
      };

      tryJoinRoom();
      
      const retryTimeout = setTimeout(() => {
        if (!hasJoinedRoom) {
          console.log('🔄 Retrying join analytics room...');
          tryJoinRoom();
          setHasJoinedRoom(true);
        }
      }, 1000);

      socket.on('analytics-update', (data) => {
        console.log('📊 Real-time analytics update:', data);
        setAnalytics(data.analytics);
        setLastUpdated(new Date());
        setHasJoinedRoom(true); 
      });

      return () => {
        clearTimeout(retryTimeout);
        socket.emit('leave-analytics', { form_id: params.id });
        socket.off('analytics-update');
        setHasJoinedRoom(false);
      };
    }
  }, [socket, params.id, hasJoinedRoom]);

  const fetchAnalytics = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await api.get(`/forms/${params.id}/analytics`);
      setAnalytics(response.data);
      setLastUpdated(new Date());
    } catch (error) {
      console.error('Error fetching analytics:', error);
      setError('Failed to load analytics');
    } finally {
      setIsLoading(false);
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

  if (error || !analytics) {
    return (
      <AuthenticatedLayout>
        <div className="text-center py-12">
          <p className="text-red-600">{error || 'Analytics not found'}</p>
          <button
            onClick={fetchAnalytics}
            className="mt-4 px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
          >
            Retry
          </button>
        </div>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-black">Analytics Dashboard</h1>
            <p className="text-gray-600">{analytics.form_title}</p>
          </div>
          <div className="flex items-center space-x-3">
            <div className="text-sm text-gray-500">
              Last updated: {lastUpdated.toLocaleTimeString()}
              <span className="ml-2 inline-block w-2 h-2 bg-green-400 rounded-full animate-pulse"></span>
            </div>
            <button
              onClick={() => router.push('/forms')}
              className="text-gray-600 hover:text-black"
            >
              ← Back to Forms
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 bg-indigo-500 rounded-md flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Total Responses</dt>
                  <dd className="text-lg font-medium text-black">{analytics.total_responses}</dd>
                </dl>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"/>
                  </svg>
                </div>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Form Fields</dt>
                  <dd className="text-lg font-medium text-black">{analytics.field_analytics.length}</dd>
                </dl>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 bg-yellow-500 rounded-md flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"/>
                  </svg>
                </div>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Completion Rate</dt>
                  <dd className="text-lg font-medium text-black">
                    {analytics.total_responses > 0 ? '100%' : '0%'}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="space-y-6">
          <h2 className="text-lg font-medium text-black">Field Analytics</h2>
          
          {analytics.field_analytics.length === 0 ? (
            <div className="bg-white rounded-lg shadow p-12 text-center">
              <div className="text-gray-400 text-6xl mb-4">📊</div>
              <h3 className="text-lg font-medium text-black mb-2">No Data Yet</h3>
              <p className="text-gray-500">
                Once people start submitting responses, analytics will appear here.
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {analytics.field_analytics.map((field) => (
                <div key={field.field_id} className="bg-white rounded-lg shadow p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h3 className="text-lg font-medium text-black">{field.field_label}</h3>
                      <p className="text-sm text-gray-500">
                        {field.field_type} • {field.response_count} responses
                      </p>
                    </div>
                    <div className="text-2xl">
                      {field.field_type === 'rating' ? '⭐' :
                       field.field_type === 'select' ? '📋' :
                       field.field_type === 'radio' ? '⚪' :
                       field.field_type === 'checkbox' ? '☑️' :
                       field.field_type === 'number' ? '🔢' : '📝'}
                    </div>
                  </div>

                  <div className="space-y-3">
                    {field.field_type === 'rating' && field.data.average_rating && (
                      <div>
                        <div className="flex items-center space-x-2 mb-2">
                          <span className="text-sm font-medium text-gray-900">Average Rating:</span>
                          <span className="text-lg font-bold text-yellow-600">
                            {field.data.average_rating.toFixed(1)}/5
                          </span>
                        </div>
                        {field.data.distribution && (
                          <div className="space-y-1">
                            {Object.entries(field.data.distribution).map(([rating, count]) => (
                              <div key={rating} className="flex items-center justify-between text-sm">
                                <span className="text-gray-900">{rating} stars</span>
                                <span className="font-medium text-gray-900">{count as number}</span>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    )}

                    {(field.field_type === 'select' || field.field_type === 'radio') && field.data.distribution && (
                      <div>
                        <div className="text-sm font-medium text-gray-900 mb-2">Answer Distribution:</div>
                        <div className="space-y-1">
                          {Object.entries(field.data.distribution).map(([option, count]) => (
                            <div key={option} className="flex items-center justify-between text-sm">
                              <span className="truncate text-gray-900">{option}</span>
                              <span className="font-medium ml-2 text-gray-900">{count as number}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {field.field_type === 'number' && field.data.average !== undefined && (
                      <div className="grid grid-cols-3 gap-4 text-center">
                        <div>
                          <div className="text-lg font-bold text-blue-600">{field.data.average.toFixed(2)}</div>
                          <div className="text-xs text-gray-900">Average</div>
                        </div>
                        <div>
                          <div className="text-lg font-bold text-green-600">{field.data.min}</div>
                          <div className="text-xs text-gray-900">Min</div>
                        </div>
                        <div>
                          <div className="text-lg font-bold text-red-600">{field.data.max}</div>
                          <div className="text-xs text-gray-900">Max</div>
                        </div>
                      </div>
                    )}

                    {(field.field_type === 'text' || field.field_type === 'textarea' || field.field_type === 'email') && (
                      <div className="text-center">
                        <div className="text-lg font-bold text-indigo-600">
                          {field.data.average_length?.toFixed(1) || 0}
                        </div>
                        <div className="text-xs text-gray-900">Avg. Characters</div>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="text-center py-4">
          <div className="inline-flex items-center space-x-2 text-sm text-gray-500">
            <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
            <span>Live updates enabled</span>
          </div>
        </div>
      </div>
    </AuthenticatedLayout>
  );
}