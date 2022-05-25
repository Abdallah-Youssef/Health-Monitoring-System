from django.shortcuts import render,redirect

from django_app.duckdpquery import query
from .forms import QueryForm
from django.http import HttpResponse
import datetime
import django_app.formater as formater

def query_result_view(request,start_date,end_date):
    start_date, end_date = formater.format_dates(start_date, end_date)
    result = query(start_date, end_date)
    # query duckdb then pass the query result as the first argument in format_result
    return HttpResponse(formater.format_result(result,start_date,end_date))


def query_form_view(request):

    # submit button is pushed
    if request.method == 'POST':
        # extract start time and end time
        start_date = request.POST['start_date']
        end_date = request.POST['end_date']
        # get the form to validate the inputs
        form = QueryForm(request.POST)
        if form.is_valid():
            # convert into datetime to make sure that end time > start time
            start_date_obj = datetime.datetime.strptime(start_date, '%d/%m/%Y %H:%M')
            end_date_obj = datetime.datetime.strptime(end_date, '%d/%m/%Y %H:%M')
            if end_date_obj < start_date_obj:
                # end time is before start time so it will preview an error msg
                context = {'form': QueryForm(), 'error_msg': "end date cannot be before start date!"}
                return render(request, "queryForm.html", context)
            else:
                #inputs are valid
                # base_url = reverse('query')
                # print(base_url)
                # query_string=urlencode({'start_date':str(start_date),'end_date':str(end_date)})
                # print(query_string)
                # url = '{}?{}'.format(base_url, query_string)
                return redirect("query",start_date=str(start_date),end_date=str(end_date))
        else:
            context = {'form': QueryForm(), 'error_msg': "follow the input format: d/m/Y H:M !"}
            return render(request, "queryForm.html", context)
    #calling the form page
    elif request.method == 'GET':
        context = {'form': QueryForm()}
        return render(request, "queryForm.html", context)


