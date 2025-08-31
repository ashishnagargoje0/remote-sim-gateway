import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { Eye, EyeOff, Smartphone, Loader2 } from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import Button from '../components/common/Button';
import Input from '../components/common/Input';

export default function Login() {
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();
  const { login, user, loading } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm();

  // Redirect if already authenticated
  useEffect(() => {
    if (!loading && user) {
      router.push('/');
    }
  }, [user, loading, router]);

  const onSubmit = async (data) => {
    setIsLoading(true);
    
    const result = await login(data.email, data.password);
    
    if (!result.success) {
      setError('root', {
        type: 'manual',
        message: result.error,
      });
    }
    
    setIsLoading(false);
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <>
      <Head>
        <title>Login - Remote SIM Gateway</title>
        <meta name="description" content="Login to Remote SIM Gateway" />
      </Head>

      <div className="min-h-screen flex">
        {/* Left side - Branding */}
        <div className="hidden lg:flex lg:flex-1 lg:flex-col lg:justify-center lg:px-12 xl:px-16 bg-gradient-to-br from-primary-600 to-primary-800">
          <div className="mx-auto max-w-md">
            <div className="flex items-center mb-8">
              <div className="h-12 w-12 bg-white rounded-xl flex items-center justify-center">
                <Smartphone className="h-8 w-8 text-primary-600" />
              </div>
              <div className="ml-4">
                <h1 className="text-2xl font-bold text-white">Remote SIM Gateway</h1>
                <p className="text-primary-100">Control your phone remotely</p>
              </div>
            </div>
            
            <div className="space-y-6 text-white">
              <div className="flex items-start space-x-3">
                <div className="flex-shrink-0 h-6 w-6 bg-primary-500 rounded-full flex items-center justify-center mt-0.5">
                  <span className="text-sm font-bold">1</span>
                </div>
                <div>
                  <h3 className="font-semibold">Send SMS Remotely</h3>
                  <p className="text-primary-100 text-sm">Send text messages from anywhere in the world using your phone's SIM card</p>
                </div>
              </div>
              
              <div className="flex items-start space-x-3">
                <div className="flex-shrink-0 h-6 w-6 bg-primary-500 rounded-full flex items-center justify-center mt-0.5">
                  <span className="text-sm font-bold">2</span>
                </div>
                <div>
                  <h3 className="font-semibold">Make Calls</h3>
                  <p className="text-primary-100 text-sm">Initiate phone calls remotely through your connected Android device</p>
                </div>
              </div>
              
              <div className="flex items-start space-x-3">
                <div className="flex-shrink-0 h-6 w-6 bg-primary-500 rounded-full flex items-center justify-center mt-0.5">
                  <span className="text-sm font-bold">3</span>
                </div>
                <div>
                  <h3 className="font-semibold">Monitor & Manage</h3>
                  <p className="text-primary-100 text-sm">Track all your communications with detailed logs and analytics</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Right side - Login form */}
        <div className="flex-1 flex flex-col justify-center px-4 sm:px-6 lg:px-8 xl:px-12">
          <div className="mx-auto w-full max-w-md">
            {/* Mobile logo */}
            <div className="lg:hidden flex items-center justify-center mb-8">
              <div className="h-12 w-12 bg-primary-600 rounded-xl flex items-center justify-center">
                <Smartphone className="h-8 w-8 text-white" />
              </div>
              <div className="ml-4">
                <h1 className="text-2xl font-bold text-gray-900">SIM Gateway</h1>
              </div>
            </div>

            <div>
              <h2 className="text-3xl font-bold text-gray-900">Sign in to your account</h2>
              <p className="mt-2 text-sm text-gray-600">
                Don't have an account?{' '}
                <Link href="/register" className="font-medium text-primary-600 hover:text-primary-500">
                  Sign up for free
                </Link>
              </p>
            </div>

            <div className="mt-8">
              <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
                {/* Email */}
                <div>
                  <Input
                    label="Email address"
                    type="email"
                    autoComplete="email"
                    required
                    {...register('email', {
                      required: 'Email is required',
                      pattern: {
                        value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                        message: 'Invalid email address',
                      },
                    })}
                    error={errors.email?.message}
                  />
                </div>

                {/* Password */}
                <div>
                  <div className="relative">
                    <Input
                      label="Password"
                      type={showPassword ? 'text' : 'password'}
                      autoComplete="current-password"
                      required
                      {...register('password', {
                        required: 'Password is required',
                        minLength: {
                          value: 6,
                          message: 'Password must be at least 6 characters',
                        },
                      })}
                      error={errors.password?.message}
                    />
                    <button
                      type="button"
                      className="absolute inset-y-0 right-0 top-6 pr-3 flex items-center"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? (
                        <EyeOff className="h-5 w-5 text-gray-400" />
                      ) : (
                        <Eye className="h-5 w-5 text-gray-400" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Remember me and Forgot password */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <input
                      id="remember-me"
                      name="remember-me"
                      type="checkbox"
                      className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                    />
                    <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-900">
                      Remember me
                    </label>
                  </div>

                  <div className="text-sm">
                    <Link href="/forgot-password" className="font-medium text-primary-600 hover:text-primary-500">
                      Forgot your password?
                    </Link>
                  </div>
                </div>

                {/* Error message */}
                {errors.root && (
                  <div className="rounded-md bg-red-50 p-4">
                    <div className="text-sm text-red-700">{errors.root.message}</div>
                  </div>
                )}

                {/* Submit button */}
                <div>
                  <Button
                    type="submit"
                    className="w-full"
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <>
                        <Loader2 className="animate-spin h-4 w-4 mr-2" />
                        Signing in...
                      </>
                    ) : (
                      'Sign in'
                    )}
                  </Button>
                </div>
              </form>

              {/* Demo credentials */}
              <div className="mt-6 p-4 bg-gray-50 rounded-lg">
                <p className="text-sm font-medium text-gray-700 mb-2">Demo Credentials:</p>
                <p className="text-xs text-gray-600">Email: demo@simgateway.com</p>
                <p className="text-xs text-gray-600">Password: demo123</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}