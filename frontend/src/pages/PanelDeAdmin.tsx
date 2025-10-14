import type { FC, FormEvent } from "react";
import { useState } from "react";
import { useNavigate } from 'react-router-dom'
import { logoutStaff } from '../lib/auth'
import Boton from "../components/Boton";
import "./PanelDeAdmin.css";

type EntityType =
	| "software"
	| "gabinete"
	| "procesador"
	| "placaGrafica"
	| "almacenamiento"
	| "placaMadre"
	| "fuente"
	| "memoriaRam";

type GenericItem = {
	id: string;
	tipo: EntityType;
	nombre: string;
	orden?: number;
	[key: string]: any;
};

const emptyTemplate = (tipo: EntityType): GenericItem => {
	const base: GenericItem = { id: "", tipo, nombre: "", orden: 0 };
	switch (tipo) {
		case "software":
			return { ...base, categoria: "", empresa: "", orden_grafica: 0, orden_procesador: 0, orden_ram: 0 };
		case "gabinete":
			return { ...base, fabricante: "", formFactor: "", imageURL: "", maxLargoGpuFloat: 0, msrp: 0, socket: "" };
		case "procesador":
			return { ...base, consumo: 0, fabricante: "", generacion: "", imageURL: "", msrp: 0, nucleos: 0 };
		case "placaGrafica":
			return { ...base, consumo: 0, fabricante: "", imageURL: "", largo: 0, modelo: "", vram: 0 };
		case "almacenamiento":
			return { ...base, capacidad: 0, fabricante: "", imageURL: "", msrp: 0, tipo_almacenamiento: "" };
		case "placaMadre":
			return { ...base, fabricante: "", formFactor: "", imageURL: "", modelo: "", msrp: 0, socketProcesador: "", socketRam: "" };
		case "fuente":
			return { ...base, capacidad: 0, fabricante: "", imageURL: "" };
		case "memoriaRam":
			return { ...base, capacidad: 0, fabricante: "", generacion: "", imageURL: "" };
	}
};

const PanelDeAdmin: FC = () => {
	const [items, setItems] = useState<GenericItem[]>([]);
	const [selectedType, setSelectedType] = useState<EntityType>('software');
	const [form, setForm] = useState<GenericItem>(() => emptyTemplate('software'));
	const [editingId, setEditingId] = useState<string | null>(null);
	const [viewMode, setViewMode] = useState<'choose' | 'form' | 'list'>('choose');

	function handleChange(key: string, value: any) {
		setForm((f) => ({ ...f, [key]: value }));
	}

	function onSubmit(e: FormEvent) {
		e.preventDefault();
		// basic validation
		if (!form.nombre || !String(form.nombre).trim()) return;

		if (editingId) {
			setItems((prev) => prev.map((it) => (it.id === editingId ? { ...form, id: editingId } : it)));
			setEditingId(null);
		} else {
			const newItem = { ...form, id: crypto.randomUUID() } as GenericItem;
			setItems((prev) => [newItem, ...prev]);
		}

		setForm(emptyTemplate(selectedType));
	}

	function onEdit(id: string) {
		const it = items.find((x) => x.id === id);
		if (!it) return;
		setSelectedType(it.tipo);
		setForm(it);
		setEditingId(id);
		setViewMode('form');
	}

	function onDelete(id: string) {
		if (!confirm("Eliminar este item?")) return;
		setItems((prev) => prev.filter((x) => x.id !== id));
		if (editingId === id) {
			setEditingId(null);
			setForm(emptyTemplate(selectedType));
		}
	}

	const navigate = useNavigate()

	function doLogout() {
		logoutStaff()
		navigate('/login')
	}

	return (
		<div className="admin-panel-root">
			<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
				<h1>Panel de Administración - Softwares</h1>
				<div>
					<button onClick={doLogout} className="btn-cancel">Cerrar sesión</button>
				</div>
			</div>

			{viewMode === 'choose' ? (
				<div className="panel-choices">
					<p>Seleccione la clase que desea gestionar</p>

					<div className="entity-grid">
						{(
							[
								{ key: 'software', label: 'Software' },
								{ key: 'gabinete', label: 'Gabinete' },
								{ key: 'procesador', label: 'Procesador' },
								{ key: 'placaGrafica', label: 'Placa Gráfica' },
								{ key: 'almacenamiento', label: 'Almacenamiento' },
								{ key: 'placaMadre', label: 'Placa Madre' },
								{ key: 'fuente', label: 'Fuente de poder' },
								{ key: 'memoriaRam', label: 'Memoria RAM' },
							] as { key: EntityType; label: string }[]
						).map((t) => (
							<div
								key={t.key}
								className={`entity-card ${selectedType === t.key ? 'active' : ''}`}
								onClick={() => { setSelectedType(t.key); setForm(emptyTemplate(t.key)); setEditingId(null); setViewMode('form'); }}
							>
								<div className="entity-card-inner" aria-hidden>
									<div className="entity-icon">🧩</div>
								</div>

								<div style={{ padding: 8 }}>
									<div className="entity-name">{t.label}</div>
									<div className="entity-actions">
										<button
											className="entity-action-btn"
											onClick={(e) => { e.stopPropagation(); setSelectedType(t.key); setForm(emptyTemplate(t.key)); setViewMode('form'); }}
										>
											Agregar
										</button>
										<button
											className="entity-action-btn"
											onClick={(e) => { e.stopPropagation(); setSelectedType(t.key); setViewMode('list'); }}
										>
											Ver
										</button>
									</div>
								</div>
							</div>
						))}
					</div>
				</div>
			) : (
			<div className="admin-panel">
				<div className="panel-top-actions">
					<button className="choice-back" onClick={() => setViewMode('choose')}>Volver</button>
				</div>

				{viewMode === 'form' && (
				<form className="admin-form" onSubmit={onSubmit}>
					<h2>{editingId ? `Editar ${form.tipo}` : `Nuevo ${selectedType}`}</h2>

					{/* common name */}
					<label>
						Nombre
						<input value={form.nombre} onChange={(e) => handleChange('nombre', e.target.value)} required />
					</label>

					{/* dynamic fields */}
					{form.tipo === 'software' || selectedType === 'software' ? (
						<>
							<label>
								Categoria
								<input value={form.categoria || ''} onChange={(e) => handleChange('categoria', e.target.value)} />
							</label>
							<label>
								Empresa
								<input value={form.empresa || ''} onChange={(e) => handleChange('empresa', e.target.value)} />
							</label>
							<label>
								Orden gráfica
								<input type="number" value={form.orden_grafica || 0} onChange={(e) => handleChange('orden_grafica', Number(e.target.value))} />
							</label>
							<label>
								Orden procesador
								<input type="number" value={form.orden_procesador || 0} onChange={(e) => handleChange('orden_procesador', Number(e.target.value))} />
							</label>
							<label>
								Orden RAM
								<input type="number" value={form.orden_ram || 0} onChange={(e) => handleChange('orden_ram', Number(e.target.value))} />
							</label>
						</>
					) : null}

					{selectedType === 'gabinete' || form.tipo === 'gabinete' ? (
						<>
							<label>
								Fabricante
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} />
							</label>
							<label>
								Form-factor
								<input value={form.formFactor || ''} onChange={(e) => handleChange('formFactor', e.target.value)} />
							</label>
							<label>
								Image URL
								<input value={form.imageURL || ''} onChange={(e) => handleChange('imageURL', e.target.value)} />
							</label>
							<label>
								Max largo GPU (float)
								<input type="number" step="0.1" value={form.maxLargoGpuFloat || 0} onChange={(e) => handleChange('maxLargoGpuFloat', Number(e.target.value))} />
							</label>
							<label>
								MSRP
								<input type="number" step="0.01" value={form.msrp || 0} onChange={(e) => handleChange('msrp', Number(e.target.value))} />
							</label>
							<label>
								Socket
								<input value={form.socket || ''} onChange={(e) => handleChange('socket', e.target.value)} />
							</label>
						</>
					) : null}

					{selectedType === 'procesador' || form.tipo === 'procesador' ? (
						<>
							<label>
								Consumo
								<input type="number" step="0.1" value={form.consumo || 0} onChange={(e) => handleChange('consumo', Number(e.target.value))} />
							</label>
							<label>
								Fabricante
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} />
							</label>
							<label>
								Generación
								<input value={form.generacion || ''} onChange={(e) => handleChange('generacion', e.target.value)} />
							</label>
							<label>
								Image URL
								<input value={form.imageURL || ''} onChange={(e) => handleChange('imageURL', e.target.value)} />
							</label>
							<label>
								MSRP
								<input type="number" step="0.01" value={form.msrp || 0} onChange={(e) => handleChange('msrp', Number(e.target.value))} />
							</label>
							<label>
								Nucleos
								<input type="number" value={form.nucleos || 0} onChange={(e) => handleChange('nucleos', Number(e.target.value))} />
							</label>
						</>
					) : null}

					{selectedType === 'placaGrafica' || form.tipo === 'placaGrafica' ? (
						<>
							<label>
								Consumo
								<input type="number" step="0.1" value={form.consumo || 0} onChange={(e) => handleChange('consumo', Number(e.target.value))} />
							</label>
							<label>
								Fabricante
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} />
							</label>
							<label>
								Image URL
								<input value={form.imageURL || ''} onChange={(e) => handleChange('imageURL', e.target.value)} />
							</label>
							<label>
								Largo
								<input type="number" step="0.1" value={form.largo || 0} onChange={(e) => handleChange('largo', Number(e.target.value))} />
							</label>
							<label>
								Modelo
								<input value={form.modelo || ''} onChange={(e) => handleChange('modelo', e.target.value)} />
							</label>
							<label>
								VRAM
								<input type="number" value={form.vram || 0} onChange={(e) => handleChange('vram', Number(e.target.value))} />
							</label>
						</>
					) : null}

					{selectedType === 'almacenamiento' || form.tipo === 'almacenamiento' ? (
						<>
							<label>
								Capacidad
								<input type="number" step="0.1" value={form.capacidad || 0} onChange={(e) => handleChange('capacidad', Number(e.target.value))} />
							</label>
							<label>
								Fabricante
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} />
							</label>
							<label>
								Image URL
								<input value={form.imageURL || ''} onChange={(e) => handleChange('imageURL', e.target.value)} />
							</label>
							<label>
								MSRP
								<input type="number" step="0.01" value={form.msrp || 0} onChange={(e) => handleChange('msrp', Number(e.target.value))} />
							</label>
							<label>
								Tipo
								<input value={form.tipo_almacenamiento || ''} onChange={(e) => handleChange('tipo_almacenamiento', e.target.value)} />
							</label>
						</>
					) : null}

					{selectedType === 'placaMadre' || form.tipo === 'placaMadre' ? (
						<>
							<label>
								Fabricante
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} />
							</label>
							<label>
								Form-factor
								<input value={form.formFactor || ''} onChange={(e) => handleChange('formFactor', e.target.value)} />
							</label>
							<label>
								Image URL
								<input value={form.imageURL || ''} onChange={(e) => handleChange('imageURL', e.target.value)} />
							</label>
							<label>
								Modelo
								<input value={form.modelo || ''} onChange={(e) => handleChange('modelo', e.target.value)} />
							</label>
							<label>
								MSRP
								<input type="number" step="0.01" value={form.msrp || 0} onChange={(e) => handleChange('msrp', Number(e.target.value))} />
							</label>
							<label>
								Socket procesador
								<input value={form.socketProcesador || ''} onChange={(e) => handleChange('socketProcesador', e.target.value)} />
							</label>
							<label>
								Socket RAM
								<input value={form.socketRam || ''} onChange={(e) => handleChange('socketRam', e.target.value)} />
							</label>
						</>
					) : null}

					{selectedType === 'fuente' || form.tipo === 'fuente' ? (
						<>
							<label>
								Capacidad
								<input type="number" step="0.1" value={form.capacidad || 0} onChange={(e) => handleChange('capacidad', Number(e.target.value))} />
							</label>
							<label>
								Fabricante
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} />
							</label>
							<label>
								Image URL
								<input value={form.imageURL || ''} onChange={(e) => handleChange('imageURL', e.target.value)} />
							</label>
						</>
					) : null}

					{selectedType === 'memoriaRam' || form.tipo === 'memoriaRam' ? (
						<>
							<label>
								Capacidad
								<input type="number" value={form.capacidad || 0} onChange={(e) => handleChange('capacidad', Number(e.target.value))} />
							</label>
							<label>
								Fabricante
								<input value={form.fabricante || ''} onChange={(e) => handleChange('fabricante', e.target.value)} />
							</label>
							<label>
								Generación
								<input value={form.generacion || ''} onChange={(e) => handleChange('generacion', e.target.value)} />
							</label>
							<label>
								Image URL
								<input value={form.imageURL || ''} onChange={(e) => handleChange('imageURL', e.target.value)} />
							</label>
						</>
					) : null}

					<div className="form-actions">
						<Boton>{editingId ? "Guardar" : "Crear"}</Boton>
						<button
							type="button"
							className="btn-cancel"
							onClick={() => {
								setForm(emptyTemplate(selectedType));
								setEditingId(null);
							}}
						>
							Cancelar
						</button>
					</div>
				</form>
				)}

				<div className="admin-list" style={{ display: viewMode === 'list' ? undefined : (viewMode === 'form' ? 'none' : undefined) }}>
					<h2>Items cargados - {selectedType} ({items.filter(i => i.tipo === selectedType).length})</h2>
					{items.filter(i => i.tipo === selectedType).length === 0 ? (
						<div className="empty">no hay nada cargado aun</div>
					) : (
						<ul>
							{items.filter(i => i.tipo === selectedType).map((it) => (
								<li key={it.id} className="list-item">
									<div className="item-main">
										<div>
											<strong>{it.nombre}</strong>
											<div className="muted">{it.fabricante || it.categoria || ''} — {it.modelo || it.empresa || ''}</div>
										</div>
										<div className="item-orders">
											{it.vram ? `VRAM:${it.vram}` : ''} {it.nucleos ? `N:${it.nucleos}` : ''} {it.capacidad ? `C:${it.capacidad}` : ''}
										</div>
									</div>

									<div className="item-actions">
										<button className="btn-edit" onClick={() => onEdit(it.id)}>Editar</button>
										<button className="btn-delete" onClick={() => onDelete(it.id)}>Borrar</button>
									</div>
								</li>
							))}
						</ul>
					)}
				</div>
			</div>
			)}
		</div>
	);
};

export default PanelDeAdmin;

