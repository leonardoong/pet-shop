import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="flex items-center justify-center py-20">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-gray-300 mb-2">404</h1>
        <p className="text-lg text-gray-500 mb-6">Page not found</p>
        <Link to="/dashboard" className="inline-flex px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700">
          Back to Dashboard
        </Link>
      </div>
    </div>
  )
}
