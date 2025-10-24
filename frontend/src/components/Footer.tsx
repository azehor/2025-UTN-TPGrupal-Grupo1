import type { FC } from 'react'

const Footer: FC = () => {
  return (
    <footer className="bg-[#101c22] border-t border-gray-700/50">
      <div className="container mx-auto px-6 py-6">
        {/* Sponsors Section */}
        <div className="mb-4">
          <h3 className="text-lg font-bold text-white mb-3 text-center">Nuestros Sponsors</h3>
          
          {/* All Sponsors in compact grid */}
          <div className="grid grid-cols-3 md:grid-cols-6 lg:grid-cols-8 gap-3 text-center text-sm">
            {/* Hardware */}
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Intel</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">AMD</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">NVIDIA</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">ASUS</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Corsair</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Samsung</div>
            </div>
            {/* Software */}
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Steam</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Adobe</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Microsoft</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Autodesk</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">JetBrains</div>
            </div>
            {/* Stores */}
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Newegg</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Amazon</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Best Buy</div>
            </div>
            <div className="text-gray-400 hover:text-white transition-colors">
              <div className="font-medium">Micro Center</div>
            </div>
          </div>
        </div>

        <div className="border-t border-gray-700/50 pt-3 text-center">
          <p className="text-gray-500 text-sm">
            © {new Date().getFullYear()} QuéPC. Proyecto académico UTN
          </p>
        </div>
      </div>
    </footer>
  )
}

export default Footer
