import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Hero section */}
      <div className="relative pt-16 pb-32 flex content-center items-center justify-center">
        <div className="container relative mx-auto">
          <div className="items-center flex flex-wrap">
            <div className="w-full lg:w-6/12 px-4 ml-auto mr-auto text-center">
              <div className="pr-12">
                <h1 className="text-5xl font-bold text-gray-900 dark:text-white">
                  Hello Pulse
                </h1>
                <p className="mt-4 text-lg text-gray-600 dark:text-gray-300">
                  Collaborative brainstorming platform powered by artificial intelligence. 
                  Transform your team's collaboration with intuitive tools and AI assistance.
                </p>
                <div className="mt-10">
                  <Link 
                    href="/register" 
                    className="inline-block px-6 py-3 mr-3 font-bold text-center text-white uppercase align-middle transition-all rounded-lg cursor-pointer bg-blue-600 hover:bg-blue-700 text-sm"
                  >
                    Get Started
                  </Link>
                  <Link 
                    href="/login" 
                    className="inline-block px-6 py-3 font-bold text-center text-gray-800 uppercase align-middle transition-all bg-white rounded-lg cursor-pointer border border-gray-300 hover:bg-gray-100 text-sm"
                  >
                    Sign In
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Features section */}
      <section className="py-12 bg-white dark:bg-gray-800">
        <div className="container mx-auto px-4">
          <div className="flex flex-wrap justify-center text-center mb-12">
            <div className="w-full lg:w-6/12 px-4">
              <h2 className="text-3xl font-semibold text-gray-900 dark:text-white">
                Key Features
              </h2>
              <p className="text-lg leading-relaxed m-4 text-gray-600 dark:text-gray-300">
                Everything you need to boost your team's collaboration and productivity.
              </p>
            </div>
          </div>
          <div className="flex flex-wrap">
            <div className="lg:w-1/3 px-4 text-center mb-10">
              <div className="h-12 w-12 text-blue-500 rounded-full flex items-center justify-center text-3xl mx-auto mb-4">
                🧠
              </div>
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                AI-Powered Insights
              </h3>
              <p className="mt-2 mb-4 text-gray-600 dark:text-gray-400">
                Get intelligent suggestions and insights from our advanced AI.
              </p>
            </div>
            <div className="lg:w-1/3 px-4 text-center mb-10">
              <div className="h-12 w-12 text-blue-500 rounded-full flex items-center justify-center text-3xl mx-auto mb-4">
                👥
              </div>
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                Team Collaboration
              </h3>
              <p className="mt-2 mb-4 text-gray-600 dark:text-gray-400">
                Work together with your team in real-time on projects and ideas.
              </p>
            </div>
            <div className="lg:w-1/3 px-4 text-center mb-10">
              <div className="h-12 w-12 text-blue-500 rounded-full flex items-center justify-center text-3xl mx-auto mb-4">
                📊
              </div>
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                Integrated Organization
              </h3>
              <p className="mt-2 mb-4 text-gray-600 dark:text-gray-400">
                Manage projects, events, and team members in one place.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-6 bg-gray-100 dark:bg-gray-900">
        <div className="container mx-auto px-4">
          <div className="flex flex-wrap items-center md:justify-between justify-center">
            <div className="w-full md:w-4/12 px-4 mx-auto text-center">
              <div className="text-sm text-gray-600 dark:text-gray-400 font-semibold py-1">
                Hello Pulse © {new Date().getFullYear()}
              </div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}