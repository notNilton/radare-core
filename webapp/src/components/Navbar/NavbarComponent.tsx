import React, { useState } from "react";
import "./NavbarComponent.scss";
import AboutModal from "../About/AboutModal"; // Importa o componente AboutModal

interface NavbarComponentProps {
  version: string;
}

const NavbarComponent: React.FC<NavbarComponentProps> = () => {
  const [showAbout, setShowAbout] = useState<boolean>(false);

  const toggleAboutPopup = () => {
    setShowAbout(!showAbout);
  };

  return (
    <>
      <header className="header">
        <nav className="navbar">
          <ul className="navbar-menu">
            <li className="navbar-item">
              <button className="link-button">Arquivo</button>
              <ul>
                <li>
                  <button className="link-button">Novo</button>
                </li>
                <li>
                  <button className="link-button">Abrir</button>
                </li>
                <li>
                  <button className="link-button">Fechar</button>
                </li>
                <li>
                  <button className="link-button">Salvar</button>
                </li>
                <li>
                  <button className="link-button">Salvar Como</button>
                </li>
                <li>
                  <button className="link-button">Importar</button>
                </li>
                <li>
                  <button className="link-button">Exportar</button>
                </li>
                <li>
                  <button className="link-button">Exemplos</button>
                </li>
                <li>
                  <button className="link-button">Sair</button>
                </li>
              </ul>
            </li>
            <li className="navbar-item">
              <button className="link-button">Editar</button>
              <ul>
                <li>
                  <button className="link-button">Retornar</button>
                </li>
                <li>
                  <button className="link-button">Avançar</button>
                </li>
                <li>
                  <button className="link-button">Cortar</button>
                </li>
                <li>
                  <button className="link-button">Copiar</button>
                </li>
                <li>
                  <button className="link-button">Colar</button>
                </li>
                <li>
                  <button className="link-button">Deletar</button>
                </li>
              </ul>
            </li>
            <li className="navbar-item">
              <button className="link-button">Desenhar</button>
              <ul>
                <li>
                  <button className="link-button">
                    Lista de Objetos Padrões
                  </button>
                </li>
                <li>
                  <button className="link-button">Agrupar Objetos</button>
                </li>
                <li>
                  <button className="link-button">Desagrupar Objetos</button>
                </li>
              </ul>
            </li>
            <li className="navbar-item">
              <button className="link-button">Ajuda</button>
            </li>
            <li className="navbar-item">
              <button className="link-button" onClick={toggleAboutPopup}>
                Sobre
              </button>
            </li>
          </ul>
        </nav>
      </header>
      <AboutModal showAbout={showAbout} toggleAboutPopup={toggleAboutPopup} />
    </>
  );
};

export default NavbarComponent;
