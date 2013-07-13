$(function() {
	configurarBusca();
})

function configurarBusca() {
	var comprimento = 3;
	var buscaNoServidor = "";
	var buscaAnterior = "";
	var resultados = [];

	var obterInscritos = function(busca) {
		$.get("adm/busca/lideres?contem=" + busca, function(lideres) {
			var lista = [];
			for (i in lideres) {
				var id = lideres[i].id
				var elemento = "<a href='/adm/lider/" + id + "'>" +
				lideres[i].nome + "<br/></a>";

				lista.push(elemento);
			}

			obterVoluntarios(lista, busca)
		});
	}

	var obterVoluntarios = function(inscritos, busca) {
		$.get("adm/busca/voluntarios?contem=" + busca, function(voluntarios) {
			var lista = [];
			for (i in voluntarios) {
				var id = voluntarios[i].id
				var elemento = "<a href='/adm/voluntario/" + id + "'>" +
				voluntarios[i].nome + "<br/></a>";

				lista.push(elemento);
			}

			inscritos = inscritos.concat(lista)
			$("#inscritos").html(inscritos);
			resultados = inscritos;
		});
	}

	$("#busca").keyup(function() {
		var novaBusca = $("#busca").val();

		if (novaBusca.length == comprimento) {
			if (novaBusca != buscaNoServidor) {
				buscaNoServidor = novaBusca;
				obterInscritos(buscaNoServidor);
			} else {
				$("#inscritos").html(resultados);
			}

		} else if (novaBusca.length < comprimento) {
			$("#inscritos").html("");

		} else {
			if (novaBusca.length < buscaAnterior.length) {
				$('#inscritos').html(resultados);
			}

			var html = $("#inscritos").children("a").filter(function() {
				re = new RegExp(novaBusca, "i");
				return $(this).text().search(re) != -1;
			})

			$('#inscritos').html(html);
		}

		buscaAnterior = novaBusca;
	});
}

