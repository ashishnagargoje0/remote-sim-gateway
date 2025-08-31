import { useState, useEffect } from 'react';
import Head from 'next/head';
import Layout from '../components/layout/Layout';
import StatsCard from '../components/dashboard/StatsCard';
import RecentActivity from '../components/dashboard/RecentActivity';
import QuickActions from '../components/dashboard/QuickActions';
import DeviceStatus from '../components/devices/DeviceStatus';
import { useAuth } from '../hooks/useAuth';
import { useWebSocket } from '../hooks/useWebSocket';
import { dashboardAPI } from '../utils/api';
import { useQuery } from '@tanstack/react-query';
import { 
  Smartphone, 
  MessageSquare, 
  Phone, 
  Activity,
  Users,
  TrendingUp
} from 'lucide-react';

export default function Dashboard() {
  const { user, loading } = useAuth();
  const { connectionStatus, connectedDevices } = useWebSocket();

  // Fetch dashboard stats
  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: dashboardAPI.getStats,
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  // Fetch recent activity
  const { data: activity, isLoading: activityLoading } = useQuery({
    queryKey: ['dashboard-activity'],
    queryFn: dashboardAPI.getRecentActivity,
    refetchInterval: 10000, // Refetch every 10 seconds
  });

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!user) {
    return null; // Will redirect to login via useAuth hook
  }

  const statsData = [
    {
      title: 'Connected Devices',
      value: connectedDevices?.length || 0,
      icon: Smartphone,
      color: 'bg-blue-500',
      change: '+2 from yesterday',
      trend: 'up'
    },
    {
      title: 'Messages Sent Today',
      value: stats?.messages_today || 0,
      icon: MessageSquare,
      color: 'bg-green-500',
      change: `+${stats?.messages_growth || 0}% from yesterday`,
      trend: stats?.messages_growth > 0 ? 'up' : 'down'
    },
    {
      title: 'Calls Made Today',
      value: stats?.calls_today || 0,
      icon: Phone,
      color: 'bg-purple-500',
      change: `+${stats?.calls_growth || 0}% from yesterday`,
      trend: stats?.calls_growth > 0 ? 'up' : 'down'
    },
    {
      title: 'Success Rate',
      value: `${stats?.success_rate || 0}%`,
      icon: TrendingUp,
      color: 'bg-emerald-500',
      change: `${stats?.success_rate_change || 0}% from yesterday`,
      trend: stats?.success_rate_change > 0 ? 'up' : 'down'
    }
  ];

  return (
    <>
      <Head>
        <title>Dashboard - Remote SIM Gateway</title>
        <meta name="description" content="Remote SIM Gateway Dashboard" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Layout>
        <div className="space-y-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
              <p className="text-gray-600 mt-1">
                Welcome back, {user.email}! Here's what's happening with your SIM gateway.
              </p>
            </div>
            
            {/* Connection Status */}
            <div className="flex items-center space-x-4">
              <div className={`flex items-center space-x-2 px-3 py-1 rounded-full text-sm font-medium ${
                connectionStatus === 'connected' 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-red-100 text-red-800'
              }`}>
                <div className={`w-2 h-2 rounded-full ${
                  connectionStatus === 'connected' ? 'bg-green-500' : 'bg-red-500'
                }`}></div>
                <span>{connectionStatus === 'connected' ? 'Connected' : 'Disconnected'}</span>
              </div>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {statsData.map((stat, index) => (
              <StatsCard key={index} {...stat} loading={statsLoading} />
            ))}
          </div>

          {/* Main Content Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Left Column - Quick Actions */}
            <div className="lg:col-span-1">
              <QuickActions />
            </div>

            {/* Middle Column - Device Status */}
            <div className="lg:col-span-1">
              <DeviceStatus devices={connectedDevices} />
            </div>

            {/* Right Column - Recent Activity */}
            <div className="lg:col-span-1">
              <RecentActivity 
                activity={activity} 
                loading={activityLoading} 
              />
            </div>
          </div>

          {/* Bottom Section - Charts/Analytics */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Message Analytics */}
            <div className="bg-white rounded-xl shadow-soft p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <MessageSquare className="w-5 h-5 mr-2 text-primary-600" />
                Message Analytics
              </h3>
              <div className="h-64">
                {/* TODO: Add chart component */}
                <div className="flex items-center justify-center h-full text-gray-500">
                  <Activity className="w-8 h-8 mr-2" />
                  Message charts will be displayed here
                </div>
              </div>
            </div>

            {/* Call Analytics */}
            <div className="bg-white rounded-xl shadow-soft p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <Phone className="w-5 h-5 mr-2 text-primary-600" />
                Call Analytics
              </h3>
              <div className="h-64">
                {/* TODO: Add chart component */}
                <div className="flex items-center justify-center h-full text-gray-500">
                  <Activity className="w-8 h-8 mr-2" />
                  Call charts will be displayed here
                </div>
              </div>
            </div>
          </div>
        </div>
      </Layout>
    </>
  );
}