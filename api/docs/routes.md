# 

## Endpoints



## Routes

<details>
<summary>`/`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/**
	- _GET_
		- [pi/internal/transport/http.(*Server).getHomepage-fm]()

</details>
<details>
<summary>`/v1/carreras`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/carreras**
		- _GET_
			- [pi/internal/carreras.(*Carreras).List-fm]()
		- _POST_
			- [pi/internal/carreras.(*Carreras).Create-fm]()

</details>
<details>
<summary>`/v1/carreras/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/carreras/{id}**
		- _DELETE_
			- [pi/internal/carreras.(*Carreras).Delete-fm]()
		- _GET_
			- [pi/internal/carreras.(*Carreras).Read-fm]()
		- _PUT_
			- [pi/internal/carreras.(*Carreras).Update-fm]()

</details>
<details>
<summary>`/v1/carreras/{id}/softwares`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/carreras/{id}/softwares**
		- _GET_
			- [pi/internal/carreraSoftware.(*CarreraSoftwares).ListByCarrera-fm]()

</details>
<details>
<summary>`/v1/carreras/{id}/softwares/{software_id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/carreras/{id}/softwares/{software_id}**
		- _POST_
			- [pi/internal/carreraSoftware.(*CarreraSoftwares).Create-fm]()
		- _DELETE_
			- [pi/internal/carreraSoftware.(*CarreraSoftwares).Delete-fm]()

</details>
<details>
<summary>`/v1/componentes/discos`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/discos**
			- _GET_
				- [pi/internal/componentes/discos.(*Discos).List-fm]()
			- _POST_
				- [pi/internal/componentes/discos.(*Discos).Create-fm]()

</details>
<details>
<summary>`/v1/componentes/discos/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/discos/{id}**
			- _DELETE_
				- [pi/internal/componentes/discos.(*Discos).Delete-fm]()
			- _GET_
				- [pi/internal/componentes/discos.(*Discos).Read-fm]()
			- _PUT_
				- [pi/internal/componentes/discos.(*Discos).Update-fm]()

</details>
<details>
<summary>`/v1/componentes/gabinetes`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/gabinetes**
			- _GET_
				- [pi/internal/componentes/gabinetes.(*Gabinetes).List-fm]()
			- _POST_
				- [pi/internal/componentes/gabinetes.(*Gabinetes).Create-fm]()

</details>
<details>
<summary>`/v1/componentes/gabinetes/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/gabinetes/{id}**
			- _DELETE_
				- [pi/internal/componentes/gabinetes.(*Gabinetes).Delete-fm]()
			- _GET_
				- [pi/internal/componentes/gabinetes.(*Gabinetes).Read-fm]()
			- _PUT_
				- [pi/internal/componentes/gabinetes.(*Gabinetes).Update-fm]()

</details>
<details>
<summary>`/v1/componentes/gpus`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/gpus**
			- _GET_
				- [pi/internal/componentes/gpus.(*GPUs).List-fm]()
			- _POST_
				- [pi/internal/componentes/gpus.(*GPUs).Create-fm]()

</details>
<details>
<summary>`/v1/componentes/gpus/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/gpus/{id}**
			- _GET_
				- [pi/internal/componentes/gpus.(*GPUs).Read-fm]()
			- _PUT_
				- [pi/internal/componentes/gpus.(*GPUs).Update-fm]()
			- _DELETE_
				- [pi/internal/componentes/gpus.(*GPUs).Delete-fm]()

</details>
<details>
<summary>`/v1/componentes/motherboards`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/motherboards**
			- _GET_
				- [pi/internal/componentes/motherboards.(*Motherboards).List-fm]()
			- _POST_
				- [pi/internal/componentes/motherboards.(*Motherboards).Create-fm]()

</details>
<details>
<summary>`/v1/componentes/motherboards/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/motherboards/{id}**
			- _PUT_
				- [pi/internal/componentes/motherboards.(*Motherboards).Update-fm]()
			- _DELETE_
				- [pi/internal/componentes/motherboards.(*Motherboards).Delete-fm]()
			- _GET_
				- [pi/internal/componentes/motherboards.(*Motherboards).Read-fm]()

</details>
<details>
<summary>`/v1/componentes/procesadores`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/procesadores**
			- _GET_
				- [pi/internal/componentes/procesadores.(*Procesadores).List-fm]()
			- _POST_
				- [pi/internal/componentes/procesadores.(*Procesadores).Create-fm]()

</details>
<details>
<summary>`/v1/componentes/procesadores/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/procesadores/{id}**
			- _PUT_
				- [pi/internal/componentes/procesadores.(*Procesadores).Update-fm]()
			- _DELETE_
				- [pi/internal/componentes/procesadores.(*Procesadores).Delete-fm]()
			- _GET_
				- [pi/internal/componentes/procesadores.(*Procesadores).Read-fm]()

</details>
<details>
<summary>`/v1/componentes/psus`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/psus**
			- _GET_
				- [pi/internal/componentes/psus.(*PSUs).List-fm]()
			- _POST_
				- [pi/internal/componentes/psus.(*PSUs).Create-fm]()

</details>
<details>
<summary>`/v1/componentes/psus/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/psus/{id}**
			- _DELETE_
				- [pi/internal/componentes/psus.(*PSUs).Delete-fm]()
			- _GET_
				- [pi/internal/componentes/psus.(*PSUs).Read-fm]()
			- _PUT_
				- [pi/internal/componentes/psus.(*PSUs).Update-fm]()

</details>
<details>
<summary>`/v1/componentes/rams`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/rams**
			- _GET_
				- [pi/internal/componentes/rams.(*RAMs).List-fm]()
			- _POST_
				- [pi/internal/componentes/rams.(*RAMs).Create-fm]()

</details>
<details>
<summary>`/v1/componentes/rams/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/componentes**
		- **/rams/{id}**
			- _PUT_
				- [pi/internal/componentes/rams.(*RAMs).Update-fm]()
			- _DELETE_
				- [pi/internal/componentes/rams.(*RAMs).Delete-fm]()
			- _GET_
				- [pi/internal/componentes/rams.(*RAMs).Read-fm]()

</details>
<details>
<summary>`/v1/recomendaciones-carrera/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/recomendaciones-carrera/{id}**
		- _GET_
			- [pi/internal/recomendaciones.(*Recomendaciones).RecomendacionCarrera-fm]()

</details>
<details>
<summary>`/v1/recomendaciones-softwares`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/recomendaciones-softwares**
		- _POST_
			- [pi/internal/recomendaciones.(*Recomendaciones).RecomendacionSoftware-fm]()

</details>
<details>
<summary>`/v1/softwares`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/softwares**
		- _GET_
			- [pi/internal/softwares.(*Softwares).List-fm]()
		- _POST_
			- [pi/internal/softwares.(*Softwares).Create-fm]()

</details>
<details>
<summary>`/v1/softwares/{id}`</summary>

- [ContentTypeMiddleware]()
- [o-chi/cors.(*Cors).Handler-fm]()
- **/v1**
	- **/softwares/{id}**
		- _GET_
			- [pi/internal/softwares.(*Softwares).Read-fm]()
		- _PUT_
			- [pi/internal/softwares.(*Softwares).Update-fm]()
		- _DELETE_
			- [pi/internal/softwares.(*Softwares).Delete-fm]()

</details>

Total # of routes: 23
