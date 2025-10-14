import type { FC } from 'react'
import './Footer.css'

const Footer: FC = () => {
  return (
    <footer className="site-footer">
      <div className="footer-inner">
        <div className="footer-section">
          <h4>QuePC</h4>
          <p>La mejor forma de armar tu pc para esa gran cursada</p>
        </div>

        <div className="footer-section">
          <h4>Redes Sociales</h4>
          <ul>
            <li><a href="#">LinkedIn - QuePC Team</a></li>
            <li><a href="#">Instagram - @quepc_official</a></li>
            <li><a href="#">X - @quepc</a></li>
            <li><a href="#">Facebook - QuePC Comunidad</a></li>
          </ul>
        </div>

        <div className="footer-section">
          <h4>Patrocinadores</h4>
          <ul>
            <li>MicroHard (falso)</li>
            <li>choogle (falso)</li>
            <li>Universidad Tecnologica Nacional (falsa)</li>
            <li>intel-LIGENTE (falso)</li>
            <li>AEME-D (falso)</li>
            <li>ENVIDIA</li>
          </ul>
        </div>
      </div>
      <div className="footer-bottom">© {new Date().getFullYear()} QuePC — Todos los derechos reservados (falso)</div>
    </footer>
  )
}

export default Footer
